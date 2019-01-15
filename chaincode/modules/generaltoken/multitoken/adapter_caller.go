package multitoken

import (
	"github.com/golang/protobuf/proto"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"hyperledger.abchain.org/chaincode/modules/generaltoken"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	txutil "hyperledger.abchain.org/core/tx"
	"math/big"
)

type GeneralCall struct {
	txgen.TxCaller
}

//we use a "simulated" caller to hook the token's caller and build our msg
type dummyCaller struct {
	*GeneralCall
}

var msgGen = map[string]func(proto.Message) *pb.MultiTokenMsg{
	generaltoken.Method_Init: func(msg proto.Message) *pb.MultiTokenMsg {
		return &pb.MultiTokenMsg{TokenMsg: &pb.MultiTokenMsg_Init{msg.(*pb.BaseToken)}}
	},
	generaltoken.Method_Transfer: func(msg proto.Message) *pb.MultiTokenMsg {
		return &pb.MultiTokenMsg{TokenMsg: &pb.MultiTokenMsg_Fund{msg.(*pb.SimpleFund)}}
	},
	generaltoken.Method_Assign: func(msg proto.Message) *pb.MultiTokenMsg {
		return &pb.MultiTokenMsg{TokenMsg: &pb.MultiTokenMsg_Fund{msg.(*pb.SimpleFund)}}
	},
	generaltoken.Method_TouchAddr: func(msg proto.Message) *pb.MultiTokenMsg {
		return &pb.MultiTokenMsg{TokenMsg: &pb.MultiTokenMsg_Query{msg.(*pb.QueryToken)}}
	},
	generaltoken.Method_QueryToken: func(msg proto.Message) *pb.MultiTokenMsg {
		return &pb.MultiTokenMsg{TokenMsg: &pb.MultiTokenMsg_Query{msg.(*pb.QueryToken)}}
	},
	generaltoken.Method_QueryGlobal: func(msg proto.Message) *pb.MultiTokenMsg {
		//queryglobal is empty so no need to pass it
		return &pb.MultiTokenMsg{}
	},
}

func (d dummyCaller) Invoke(method string, msg proto.Message) error {

}

func (d dummyCaller) Query(method string, msg proto.Message) (chan txgen.QueryResp, error) {

}
