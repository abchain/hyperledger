package blockchain

import (
	"fmt"
	"github.com/gocraft/web"
	log "github.com/op/go-logging"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/applications/util"
	"hyperledger.abchain.org/client"
	"hyperledger.abchain.org/core/tx"
	pbempty "github.com/golang/protobuf/ptypes/empty"
	"net/http"
)

var logger = log.MustGetLogger("server/blockchain")

type FabricBlockChain struct {
	*util.FabricClientBase
	cli client.ChainInfo
}

type BlockChainRouter struct {
	*web.Router
}

func CreateBlocChainRouter(root *web.Router, path string) BlockChainRouter {
	return BlockChainRouter{
		root.Subrouter(FabricBlockChain{}, path),
	}
}

const (
	FabricProxy_BlockHeight   = "height"
	FabricProxy_TransactionID = "txID"
)

func (r BlockChainRouter) Init(cfg util.FabricRPCCfg, client.ChainClient) BlockChainRouter {
	r.Middleware(func(s *FabricBlockChain, rw web.ResponseWriter,
		req *web.Request, next web.NextMiddlewareFunc) {

		c, err := cfg.GetCaller()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		s.cli, err = client.ViaRpc(c)
		if err != nil {
			http.Error(rw, "No support of fetching blockchain info via rpc", http.StatusInternalServerError)
			return
		}

		next(rw, req)
	})

	return r
}

func (r BlockChainRouter) BuildRoutes() {
	r.Get("/", (*FabricBlockChain).GetBlockchainInfo)
	r.Get("/blocks/:"+FabricProxy_BlockHeight, (*FabricBlockChain).GetBlock)
	r.Get("/transactions/:"+FabricProxy_TransactionID, (*FabricBlockChain).GetTransaction)
}

var notHyperledgerTx = `Not a hyperledger project compatible transaction`
var noParser = `No parser can be found for this transaction/event`

func (s *FabricBlockChain) handleTransaction(tx *client.ChainTransaction) *ChainTransaction {

	ret := &ChainTransaction{tx, "", "", nil, nil}

	parser, err := abchainTx.ParseTx(new(pbempty.Empty), tx.Method, tx.TxArgs)
	if err != nil {
		ret.Detail = notHyperledgerTx
		return ret
	}
	ret.Nonce = fmt.Sprintf("%X", parser.GetNounce())
	ret.ChaincodeModule = parser.GetCCname()

	if addParser, ok := registryParsers[strings.Join([]string{ret.Method, ret.ChaincodeModule}, "@")]; ok {
		//a hack: the message is always in args[1]
		msg := addParser.Msg()
		err = proto.Unmarshal(args[1], msg)
		if err != nil {
			ret.Detail = fmt.Sprintf("Invalid message arguments (%s)", err)
			return ret
		}
		ret.Data = msg
		ret.Detail = addParser.Detail(msg)
	} else {
		ret.Detail = noParser
	}
	return ret

}

func (s *FabricBlockChain) handleTxEvent(txe *client.ChainTxEvents) *ChainTxEvents {

	if addParser, ok := i.regParser[strings.Join([]string{txe.Name}, "@")]; ok {
		//a hack: the message is always in args[2]
		msg := addParser.Msg()
		err := proto.Unmarshal(txe.GetPayload(), msg)
		if err != nil {
			ret.Detail = fmt.Sprintf("Invalid event payload (%s)", err)
			return ret
		}
		ret.Detail = addParser.Detail(msg)
	} else {
		ret.Detail = noParser
	}

}

//dummy imply for parser
func (s *FabricBlockChain) HandleBlockchainInfo(i *protos.BlockchainInfo) interface{} { return i }
func (s *FabricBlockChain) HandleBlock(i *protos.Block) interface{}                   { return i }
func (s *FabricBlockChain) HandleTransaction(i *protos.Transaction) interface{}       { return i }

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
