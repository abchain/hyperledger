package multitoken

import (
	"errors"
	"github.com/golang/protobuf/proto"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"hyperledger.abchain.org/chaincode/modules/generaltoken"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"math/big"
)

type GeneralCall struct {
	txgen.TxCaller
}

//we use a "simulated" caller to hook the token's caller and build our msg
type dummyCaller struct {
	name string
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

	if g, ok := msgGen[method]; !ok {
		return errors.New("Unknown method")
	} else {

		wmsg := g(msg)
		wmsg.TokenName = d.name

		return d.TxCaller.Invoke(method, wmsg)
	}

}

func (d dummyCaller) Query(method string, msg proto.Message) (chan txgen.QueryResp, error) {

	if g, ok := msgGen[method]; !ok {
		return nil, errors.New("Unknown method")
	} else {

		wmsg := g(msg)
		wmsg.TokenName = d.name

		return d.TxCaller.Query(method, wmsg)
	}
}

func (i *GeneralCall) GetToken(name string) (generaltoken.TokenTxCore, error) {

	if err := baseNameVerifier(name); err != nil {
		return nil, err
	}

	return &generaltoken.GeneralCall{dummyCaller{name, i}}, nil

}

func (i *GeneralCall) CreateToken(name string, amount *big.Int) error {

	if err := baseNameVerifier(name); err != nil {
		return err
	}

	subg := generaltoken.GeneralCall{dummyCaller{name, i}}

	return subg.Init(amount)

}
