package nonce

import (
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/lib/caller"
	ccpb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
	"math/big"
)

type nonceQueryHandler struct {
	msg ccpb.QueryTransfer
	NonceConfig
}

type nonceAddHandler struct {
	msg ccpb.NonceData
	NonceConfig
}

func NonceQueryHandler(cfg NonceConfig) *nonceQueryHandler {
	return &nonceQueryHandler{NonceConfig: cfg}
}

func NonceAddHandler(cfg NonceConfig) *nonceAddHandler {
	return &nonceAddHandler{NonceConfig: cfg}
}

func (h *nonceQueryHandler) Msg() proto.Message { return &h.msg }
func (h *nonceAddHandler) Msg() proto.Message   { return &h.msg }

func (h *nonceQueryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := &h.msg

	err, data := h.NewTx(stub).Nonce(msg.Nonce)
	if err != nil {
		return nil, err
	}

	return rpc.EncodeRPCResult(data)
}

func (h *nonceAddHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := &h.msg

	if err := h.NewTx(stub).Add(msg.GetNoncekey(), big.NewInt(0).SetBytes(msg.GetAmount()),
		msg.GetFromLast(), msg.GetToLast()); err != nil {
		return nil, err
	}

	return []byte("done"), nil

}
