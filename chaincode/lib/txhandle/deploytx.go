package tx

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/shim"
	pb "hyperledger.abchain.org/protos"
	txutil "hyperledger.abchain.org/tx"
)

//deployTx integrating mutiple init tx handlers and dispatching them

type deployTx struct {
	msg      pb.DeployTx
	Handlers map[string]TxHandler
}

func DeployTxHandler(hs map[string]TxHandler) *deployTx {
	return &deployTx{Handlers: hs}
}

func (h *deployTx) Msg() proto.Message { return &h.msg }

func (h *deployTx) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	params := h.msg.InitParams

	for k, hh := range h.Handlers {
		arg, ok := params[k]
		if !ok {
			return nil, fmt.Errorf("No correspoinding arguments for deploying %s", k)
		}

		initmsg := hh.Msg()
		if err := proto.Unmarshal(arg, initmsg); err != nil {
			return nil, err
		}

		//not care about init's return
		if _, err := hh.Call(stub, parser); err != nil {
			return nil, fmt.Errorf("Deploy fail for handling %s: %s", k, err)
		}
	}

	return []byte("Done"), nil

}
