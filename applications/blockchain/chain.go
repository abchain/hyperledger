package blockchain

import (
	"github.com/gocraft/web"
	log "github.com/op/go-logging"
	"hyperledger.abchain.org/applications/util"
	"hyperledger.abchain.org/client"
)

var logger = log.MustGetLogger("server/blockchain")

type FabricBlockChain struct {
	*util.FabricClientBase
	cli client.ChainClient
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
