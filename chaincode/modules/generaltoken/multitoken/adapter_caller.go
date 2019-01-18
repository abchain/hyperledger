package multitoken

import (
	"github.com/golang/protobuf/proto"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"hyperledger.abchain.org/chaincode/modules/generaltoken"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
)

const (
	Method_Init        = "M" + generaltoken.Method_Init
	Method_Transfer    = "M" + generaltoken.Method_Transfer
	Method_Assign      = "M" + generaltoken.Method_Assign
	Method_TouchAddr   = "M" + generaltoken.Method_TouchAddr
	Method_QueryToken  = "M" + generaltoken.Method_QueryToken
	Method_QueryGlobal = "M" + generaltoken.Method_QueryGlobal
)

type GeneralCall struct {
	txgen.TxCaller
}

//we use a "simulated" caller to hook the token's caller and build our msg
type dummyCaller struct {
	name string
	*GeneralCall
}

func (d dummyCaller) Invoke(method string, msg proto.Message) error {

	wrapmsg, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	return d.TxCaller.Invoke("M"+method, &pb.MultiTokenMsg{TokenName: d.name, TokenMsg: wrapmsg})

}

func (d dummyCaller) Query(method string, msg proto.Message) (chan txgen.QueryResp, error) {

	wrapmsg, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return d.TxCaller.Query("M"+method, &pb.MultiTokenMsg{TokenName: d.name, TokenMsg: wrapmsg})

}

func (i *GeneralCall) GetToken(name string) (generaltoken.TokenTx, error) {

	if err := baseNameVerifier(name); err != nil {
		return nil, err
	}

	return &generaltoken.GeneralCall{dummyCaller{name, i}}, nil

}
