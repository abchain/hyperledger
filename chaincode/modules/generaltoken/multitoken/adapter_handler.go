package multitoken

import (
	"github.com/golang/protobuf/proto"
	txh "hyperledger.abchain.org/chaincode/lib/txhandle"
	"hyperledger.abchain.org/chaincode/modules/generaltoken"
	ccpb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
)

type basehandler struct {
	msg        ccpb.MultiTokenMsg
	inner      txh.TxHandler
	prehandled error
	TokenConfig
}

type getTokenError struct {
	e error
}

type noError struct{}

func (noError) Error() string { return "No error" }

func (h *basehandler) Msg() proto.Message { return &h.msg }

func (h *basehandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) (b []byte, e error) {
	defer func() {
		h.prehandled = nil
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

	h.innerPrehandled()

	if _, ok := h.prehandled.(noError); !ok {
		return nil, h.prehandled
	}

	return h.inner.Call(stub, parser)
}

func (h *basehandler) NewTx(stub shim.ChaincodeStubInterface, nc []byte) generaltoken.TokenTx {
	r, err := h.TokenConfig.NewTx(stub, nc).GetToken(h.msg.GetTokenName())
	if err != nil {
		panic(getTokenError{err})
	}

	return r
}

func (h *basehandler) innerPrehandled() {
	if h.prehandled != nil {
		return
	}

	err := proto.Unmarshal(h.msg.GetTokenMsg(), h.inner.Msg())
	if err != nil {
		h.prehandled = err
	} else {
		h.prehandled = noError{}
	}
}

func TransferHandler(cfg TokenConfig) *basehandler {
	bh := &basehandler{TokenConfig: cfg}
	bh.inner = generaltoken.TransferHandler(bh)
	return bh
}

func AssignHandler(cfg TokenConfig) *basehandler {
	bh := &basehandler{TokenConfig: cfg}
	bh.inner = generaltoken.AssignHandler(bh)
	return bh

}

func TouchHandler(cfg TokenConfig) *basehandler {
	bh := &basehandler{TokenConfig: cfg}
	bh.inner = generaltoken.TouchHandler()
	return bh
}

func TokenQueryHandler(cfg TokenConfig) *basehandler {
	bh := &basehandler{TokenConfig: cfg}
	bh.inner = generaltoken.TokenQueryHandler(bh)
	return bh
}
func GlobalQueryHandler(cfg TokenConfig) *basehandler {
	bh := &basehandler{TokenConfig: cfg}
	bh.inner = generaltoken.GlobalQueryHandler(bh)
	return bh
}

func InitHandler(cfg TokenConfig) *basehandler {
	bh := &basehandler{TokenConfig: cfg}
	bh.inner = generaltoken.InitHandler(bh)
	return bh
}
