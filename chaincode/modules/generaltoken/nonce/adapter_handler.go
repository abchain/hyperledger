package nonce

import (
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/lib/caller"
	ccpb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
	"math/big"
)

type nonceQueryHandler struct{ NonceConfig }
type nonceAddHandler struct{ NonceConfig }

func NonceQueryHandler(cfg NonceConfig) nonceQueryHandler {
	return nonceQueryHandler{NonceConfig: cfg}
}

func NonceAddHandler(cfg NonceConfig) nonceAddHandler {
	return nonceAddHandler{NonceConfig: cfg}
}

func (h nonceQueryHandler) Msg() proto.Message { return new(ccpb.QueryTransfer) }
func (h nonceAddHandler) Msg() proto.Message   { return new(ccpb.NonceData) }

func (h nonceQueryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := parser.GetMessage().(*ccpb.QueryTransfer)

	err, data := h.NewTx(stub, parser.GetNounce()).Nonce(msg.Nonce)
	if err != nil {
		return nil, err
	}

	return rpc.EncodeRPCResult(data.ToPB())
}

func (h nonceAddHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := parser.GetMessage().(*ccpb.NonceData)

	if err := h.NewTx(stub, parser.GetNounce()).Add(msg.GetNoncekey(), big.NewInt(0).SetBytes(msg.GetAmount()),
		msg.GetFromLast(), msg.GetToLast()); err != nil {
		return nil, err
	}

	return []byte("done"), nil

}
