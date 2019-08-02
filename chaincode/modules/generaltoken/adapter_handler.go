package generaltoken

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txh "hyperledger.abchain.org/chaincode/lib/txhandle"
	ccpb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
)

type transferHandler struct {
	TokenConfig
	legacy bool
}
type assignHandler struct{ TokenConfig }
type tokenQueryHandler struct{ TokenConfig }
type touchHandler struct{}
type globalQueryHandler struct{ TokenConfig }
type initHandler struct{ TokenConfig }

func TransferHandler(cfg TokenConfig) txh.TxHandler {
	return transferHandler{TokenConfig: cfg, legacy: true}
}

func TransferHandler2(cfg TokenConfig) txh.TxHandler {
	return transferHandler{TokenConfig: cfg}
}

func AssignHandler(cfg TokenConfig) txh.TxHandler {
	return assignHandler{TokenConfig: cfg}
}

func TouchHandler(TokenConfig) txh.TxHandler {
	return touchHandler{}
}

func TokenQueryHandler(cfg TokenConfig) txh.TxHandler {
	return tokenQueryHandler{TokenConfig: cfg}
}

func GlobalQueryHandler(cfg TokenConfig) txh.TxHandler {
	return globalQueryHandler{TokenConfig: cfg}
}

func InitHandler(cfg TokenConfig) txh.TxHandler {
	return initHandler{TokenConfig: cfg}
}

func (h transferHandler) Msg() proto.Message    { return new(ccpb.SimpleFund) }
func (h assignHandler) Msg() proto.Message      { return new(ccpb.SimpleFund) }
func (h touchHandler) Msg() proto.Message       { return new(ccpb.QueryToken) }
func (h tokenQueryHandler) Msg() proto.Message  { return new(ccpb.QueryToken) }
func (h globalQueryHandler) Msg() proto.Message { return new(empty.Empty) }
func (h initHandler) Msg() proto.Message        { return new(ccpb.BaseToken) }

func (h transferHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := parser.GetMessage().(*ccpb.SimpleFund)
	addrFrom, err := txutil.NewAddressFromPBMessage(msg.From)
	if err != nil {
		return nil, err
	}

	addrTo, err := txutil.NewAddressFromPBMessage(msg.To)
	if err != nil {
		return nil, err
	}

	if h.legacy {
		return h.NewTx(stub, parser.GetNonce()).Transfer(addrFrom.Hash, addrTo.Hash, toAmount(msg.Amount))
	} else {
		return h.NewTx(stub, parser.GetNonce()).Transfer2(addrFrom.Hash, addrTo.Hash, toAmount(msg.Amount))
	}

}

func (h assignHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := parser.GetMessage().(*ccpb.SimpleFund)

	addrTo, err := txutil.NewAddressFromPBMessage(msg.To)
	if err != nil {
		return nil, err
	}

	return h.NewTx(stub, parser.GetNonce()).Assign(addrTo.Hash, toAmount(msg.Amount))
}

func (h touchHandler) Call(shim.ChaincodeStubInterface, txutil.Parser) ([]byte, error) {
	return []byte("Done"), nil
}

func (h tokenQueryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := parser.GetMessage().(*ccpb.QueryToken)

	addr, err := txutil.NewAddressFromPBMessage(msg.Addr)
	if err != nil {
		return nil, err
	}

	err, data := h.NewTx(stub, parser.GetNonce()).Account(addr.Hash)
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

func (h globalQueryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	err, data := h.NewTx(stub, parser.GetNonce()).Global()
	if err != nil {
		return nil, err
	}

	return rpc.EncodeRPCResult(data.ToPB())
}

func (h initHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := parser.GetMessage().(*ccpb.BaseToken)

	token := h.NewTx(stub, parser.GetNonce())

	if err := token.Init(toAmount(msg.TotalTokens)); err != nil {
		return nil, err
	}

	return []byte("Ok"), nil

}
