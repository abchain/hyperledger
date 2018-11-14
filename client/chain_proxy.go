package client

import (
	"encoding/json"
	"fmt"
	protos "github.com/abchain/fabric/protos"
	"github.com/gocraft/web"
	"github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/spf13/viper"
	"io/ioutil"
	"math/big"
	"net/http"
)

type ChainTransaction struct {
	Height                         int64 `json:",string"`
	TxID, Chaincode, Method, Nonce string
	CreatedFlag                    bool
	//Data for the original protobuf input (Message part) and Detail left for parser
	Detail, Data interface{} `json:",omitempty"`
}

const (
	TxStatus_Success = 0
	TxStatus_Fail    = 1
)

type ChainTxEvents struct {
	TxID   string
	Status int
	Detail interface{} `json:",omitempty"`
}

type ChainBlock struct {
	Height       int64 `json:",string"`
	Hash         string
	TimeStamp    string
	Transactions []*ChainTransaction
	TxEvents     []*ChainTxEvents
}

//a parser which can handle the arguments of a transaction with purposed format in hyperledger project
type TxArgParser interface {
	Msg() proto.Message
	Detail(proto.Message) interface{}
}

type ChainClient interface {
	GetBlock(int64) *ChainBlock
	GetTransaction(string) *ChainTransaction
	GetTxEvent(string) *ChainTxEvents
	//TODO, add more methods like get range tx and filter ...

	//Registry parser
	RegParser(string, string, TxArgParser)
}

var ChainProxyViaRPC_Impls map[string]func(rpc.Caller) ChainClient
var ChainProxy_Impls map[string]func(*viper.Viper) ChainClient

//need to return JSON-tagged struct
type BlockChainParser interface {
	HandleBlockchainInfo(*protos.BlockchainInfo) interface{}
	HandleBlock(*protos.Block) interface{}
	HandleTransaction(*protos.Transaction) interface{}
}

type FabricProxy struct {
	//*FabricClientBase
	BlockChainParser
	server string
}

type FabricProxyRouter struct {
	*web.Router
}

func CreateFabricProxyRouter(root *web.Router, path string) FabricProxyRouter {
	return FabricProxyRouter{
		root.Subrouter(FabricProxy{}, path),
	}
}

const (
	FabricProxy_BlockHeight   = "height"
	FabricProxy_TransactionID = "txID"
)

func (r FabricProxyRouter) Init(server string, parser BlockChainParser) FabricProxyRouter {
	r.Middleware(func(s *FabricProxy, rw web.ResponseWriter,
		req *web.Request, next web.NextMiddlewareFunc) {

		if server == "" {
			http.Error(rw, "REST endpoint not set", http.StatusInternalServerError)
			return
		}

		s.server = server

		if parser == nil {
			s.BlockChainParser = s
		} else {
			s.BlockChainParser = parser
		}

		next(rw, req)
	})

	return r
}

func (r FabricProxyRouter) BuildRoutes() {
	r.Get("/", (*FabricProxy).GetBlockchainInfo)
	r.Get("/blocks/:"+FabricProxy_BlockHeight, (*FabricProxy).GetBlock)
	r.Get("/transactions/:"+FabricProxy_TransactionID, (*FabricProxy).GetTransaction)
}

//dummy imply for parser
func (s *FabricProxy) HandleBlockchainInfo(i *protos.BlockchainInfo) interface{} { return i }
func (s *FabricProxy) HandleBlock(i *protos.Block) interface{}                   { return i }
func (s *FabricProxy) HandleTransaction(i *protos.Transaction) interface{}       { return i }

type restError struct {
	Error string `json:"Error,omitempty"`
}

func queryBlockchain(url string, out interface{}) error {
	// HTTP Request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Requset failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Read response failed: %v", err)
	}
	//logger.Debugf("Body: %v", string(body))

	// Parse error
	errMsg := &restError{}
	json.Unmarshal(body, errMsg)
	if errMsg.Error != "" {
		return fmt.Errorf("Request failed: %v", errMsg.Error)
	}

	// Unmarshal
	err = json.Unmarshal(body, out)
	if err != nil {
		return fmt.Errorf("Unmarshal response failed: %v", err)
	}

	return nil
}

func (s *FabricProxy) GetBlockchainInfo(rw web.ResponseWriter, req *web.Request) {

	info1 := &protos.BlockchainInfo{}
	err := queryBlockchain(fmt.Sprintf("http://%s/chain", s.server), info1)
	if err != nil {
		s.NormalError(rw, err)
		return
	}
	s.Normal(rw, s.BlockChainParser.HandleBlockchainInfo(info1))
}

func (s *FabricProxy) GetBlock(rw web.ResponseWriter, req *web.Request) {

	index, ok := big.NewInt(0).SetString(req.PathParams[FabricProxy_BlockHeight], 0)
	if !ok {
		s.NormalErrorF(rw, -100, "Invalid Block count")
		return
	}

	block1 := &protos.Block{}
	err := queryBlockchain(fmt.Sprintf("http://%s/chain/blocks/%d",
		s.server, index.Int64()), block1)
	if err != nil {
		s.NormalError(rw, err)
		return
	}
	s.Normal(rw, s.BlockChainParser.HandleBlock(block1))

}

func (s *FabricProxy) GetTransaction(rw web.ResponseWriter, req *web.Request) {

	transactionID := req.PathParams[FabricProxy_TransactionID]

	if transactionID == "" {
		s.NormalErrorF(rw, -100, "Invalid TxID string")
		return
	}

	tx1 := &protos.Transaction{}
	err := queryBlockchain(fmt.Sprintf("http://%s/transactions/%s",
		s.server, transactionID), tx1)
	if err != nil {
		s.NormalError(rw, err)
		return
	}
	s.Normal(rw, s.BlockChainParser.HandleTransaction(tx1))
}
