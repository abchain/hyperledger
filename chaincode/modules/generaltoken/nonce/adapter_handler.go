package nonce

import (
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/lib/caller"
	ccpb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
)

type nonceQueryHandler struct {
	msg ccpb.QueryTransfer
	NonceConfig
}

func NonceQueryHandler(cfg NonceConfig) *nonceQueryHandler {
	return &nonceQueryHandler{NonceConfig: cfg}
}

func (h *nonceQueryHandler) Msg() proto.Message { return &h.msg }

func (h *nonceQueryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := &h.msg

	err, data := h.NewTx(stub).Nonce(msg.Nonce)
	if err != nil {
		return nil, err
	}

	return rpc.EncodeRPCResult(data)
}
