package nonce

import (
	"github.com/abchain/fabric/core/chaincode/shim"
	"github.com/golang/protobuf/proto"
	ccpb "hyperledger.abchain.org/chaincode/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txutil "hyperledger.abchain.org/tx"
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
