package generaltoken

import (
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/lib/caller"
	ccpb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
)

type FundMsg struct {
	msg ccpb.SimpleFund
}

func (m *FundMsg) Msg() *ccpb.SimpleFund {
	return &m.msg
}

type transferHandler struct {
	FundMsg
	TokenConfig
}

type assignHandler struct {
	FundMsg
	TokenConfig
}
type tokenQueryHandler struct {
	msg ccpb.QueryToken
	TokenConfig
}

type globalQueryHandler struct {
	FundMsg
	TokenConfig
}

type initHandler struct {
	msg ccpb.BaseToken
	TokenConfig
}

func TransferHandler(cfg TokenConfig) *transferHandler {
	return &transferHandler{TokenConfig: cfg}
}

func AssignHandler(cfg TokenConfig) *assignHandler {
	return &assignHandler{TokenConfig: cfg}
}

func TokenQueryHandler(cfg TokenConfig) *tokenQueryHandler {
	return &tokenQueryHandler{TokenConfig: cfg}
}
func GlobalQueryHandler(cfg TokenConfig) *globalQueryHandler {
	return &globalQueryHandler{TokenConfig: cfg}
}

func InitHandler(cfg TokenConfig) *initHandler {
	return &initHandler{TokenConfig: cfg}
}

func (h *transferHandler) Msg() proto.Message    { return &h.msg }
func (h *assignHandler) Msg() proto.Message      { return &h.msg }
func (h *tokenQueryHandler) Msg() proto.Message  { return &h.msg }
func (h *globalQueryHandler) Msg() proto.Message { return &h.msg }
func (h *initHandler) Msg() proto.Message        { return &h.msg }

func (h *transferHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := &h.msg
	addrFrom, err := txutil.NewAddressFromPBMessage(msg.From)
	if err != nil {
		return nil, err
	}

	addrTo, err := txutil.NewAddressFromPBMessage(msg.To)
	if err != nil {
		return nil, err
	}

	return h.NewTx(stub, parser.GetNounce()).Transfer(addrFrom.Hash, addrTo.Hash, toAmount(msg.Amount))
}

func (h *assignHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := &h.msg

	addrTo, err := txutil.NewAddressFromPBMessage(msg.To)
	if err != nil {
		return nil, err
	}

	return h.NewTx(stub, parser.GetNounce()).Assign(addrTo.Hash, toAmount(msg.Amount))
}

func (h *tokenQueryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := &h.msg

	addr, err := txutil.NewAddressFromPBMessage(msg.Addr)
	if err != nil {
		return nil, err
	}

	err, data := h.NewTx(stub, parser.GetNounce()).Account(addr.Hash)
	if err != nil {
		return nil, err
	}

	switch msg.Format {
	case ccpb.QueryToken_NUMBER:
		return []byte(data.Balance.Text(0)), nil
	case ccpb.QueryToken_ENCODED:
		return rpc.EncodeRPCResult(data.ToPB())
	default:
		return rpc.EncodeRPCResult(data.ToPB())
	}
}

func (h *globalQueryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	err, data := h.NewTx(stub, parser.GetNounce()).Global()
	if err != nil {
		return nil, err
	}

	return rpc.EncodeRPCResult(data.ToPB())
}

func (h *initHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := &h.msg

	token := h.NewTx(stub, parser.GetNounce())

	if err := token.Init(toAmount(msg.TotalTokens)); err != nil {
		return nil, err
	}

	return []byte("Ok"), nil

}
