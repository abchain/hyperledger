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
	msg ccpb.MultiTokenMsg
	TokenConfig
}

type GetTokenError struct {
	e error
}

func (h *basehandler) Msg() proto.Message { return &h.msg }
func (h *basehandler) NewTx(stub shim.ChaincodeStubInterface, nc []byte) generaltoken.TokenTx {
	r, err := h.TokenConfig.NewTx(stub, nc).GetToken(h.msg.GetTokenName())
	if err != nil {
		panic(GetTokenError{err})
	}

	return r
}

func (h *basehandler) innerCall(ih txh.TxHandler, stub shim.ChaincodeStubInterface, parser txutil.Parser) (r []byte, e error) {

	defer func() {
		r := recover()
		if err, ok := r.(GetTokenError); ok {
			e = err.e
			return
		} else {
			panic(r)
		}
	}()

	err := proto.Unmarshal(h.msg.GetTokenMsg(), ih.Msg())
	if err != nil {
		return nil, err
	}

	return ih.Call(stub, parser)
}

type transferHandler struct{ *basehandler }
type assignHandler struct{ *basehandler }
type tokenQueryHandler struct{ *basehandler }
type touchHandler struct{ *basehandler }
type globalQueryHandler struct{ *basehandler }
type initHandler struct{ *basehandler }

func TransferHandler(cfg TokenConfig) *transferHandler {
	return &transferHandler{&basehandler{TokenConfig: cfg}}
}

func AssignHandler(cfg TokenConfig) *assignHandler {
	return &assignHandler{&basehandler{TokenConfig: cfg}}
}

func TouchHandler() *touchHandler {
	return &touchHandler{&basehandler{}}
}

func TokenQueryHandler(cfg TokenConfig) *tokenQueryHandler {
	return &tokenQueryHandler{&basehandler{TokenConfig: cfg}}
}
func GlobalQueryHandler(cfg TokenConfig) *globalQueryHandler {
	return &globalQueryHandler{&basehandler{TokenConfig: cfg}}
}

func InitHandler(cfg TokenConfig) *initHandler {
	return &initHandler{&basehandler{TokenConfig: cfg}}
}

func (h *transferHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	return h.innerCall(generaltoken.TransferHandler(h), stub, parser)
}

func (h *assignHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	return h.innerCall(generaltoken.AssignHandler(h), stub, parser)
}

func (h *touchHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	return h.innerCall(generaltoken.TouchHandler(), stub, parser)
}

func (h *tokenQueryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	return h.innerCall(generaltoken.TokenQueryHandler(h), stub, parser)
}

func (h *globalQueryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	return h.innerCall(generaltoken.GlobalQueryHandler(h), stub, parser)
}

func (h *initHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	return h.innerCall(generaltoken.InitHandler(h), stub, parser)
}
