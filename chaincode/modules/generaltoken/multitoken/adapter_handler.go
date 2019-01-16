package multitoken

import (
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/lib/caller"
	ccpb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
)

type basehandler struct {
	msg ccpb.MultiTokenMsg
	TokenConfig
}

func (h *basehandler) Msg() proto.Message { return &h.msg }
