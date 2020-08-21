package multitoken

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	txh "hyperledger.abchain.org/chaincode/lib/txhandle"
	"hyperledger.abchain.org/chaincode/modules/generaltoken"
	ccpb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
)

type basehandler struct {
	innerH func(generaltoken.TokenConfig) txh.TxHandler
	TokenConfig
}

func (h basehandler) Msg() proto.Message { return new(ccpb.MultiTokenMsg) }

func (h basehandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) (rbt []byte, e error) {

	msg := parser.GetMessage().(*ccpb.MultiTokenMsg)
	defer parser.PopMsg()

	switch m := msg.GetMsg().(type) {
	case *ccpb.MultiTokenMsg_Fund:
		parser.PushMsg(m.Fund)
	case *ccpb.MultiTokenMsg_Query:
		parser.PushMsg(m.Query)
	case *ccpb.MultiTokenMsg_Init:
		parser.PushMsg(m.Init)
	default:
		parser.PushMsg(&empty.Empty{})
	}

	defer func() {
		r := recover()
		if r == nil {
			return
		} else if err, ok := r.(getTokenError); ok {
			e = err.e
			return
		} else {
			panic(r)
		}
	}()

	return h.innerH(subhandler{msg.GetTokenName(), h.TokenConfig}).Call(stub, parser)
}

type getTokenError struct {
	e error
}

type subhandler struct {
	name string
	TokenConfig
}

func (h subhandler) NewTx(stub shim.ChaincodeStubInterface, nc []byte) generaltoken.TokenTx {
	r, err := h.TokenConfig.NewTx(stub, nc).GetToken(h.name)
	if err != nil {
		panic(err)
	}

	return r
}

func TransferHandler(cfg TokenConfig) txh.TxHandler {
	return basehandler{generaltoken.TransferHandler, cfg}
}

func AssignHandler(cfg TokenConfig) txh.TxHandler {
	return basehandler{generaltoken.AssignHandler, cfg}
}

func TouchHandler(cfg TokenConfig) txh.TxHandler {
	return basehandler{generaltoken.TouchHandler, cfg}
}

func TokenQueryHandler(cfg TokenConfig) txh.TxHandler {
	return basehandler{generaltoken.TokenQueryHandler, cfg}
}
func GlobalQueryHandler(cfg TokenConfig) txh.TxHandler {
	return basehandler{generaltoken.GlobalQueryHandler, cfg}
}

func InitHandler(cfg TokenConfig) txh.TxHandler {
	return basehandler{generaltoken.InitHandler, cfg}
}
