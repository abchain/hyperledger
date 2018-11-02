package tx

import (
	"github.com/golang/protobuf/proto"
	pb "hyperledger.abchain.org/protos"
)

type DeployTxCall struct {
	InitParams map[string]proto.Message
	*TxGenerator
}

func NewDeployTx() *DeployTxCall {

	ret := new(DeployTxCall)
	ret.InitParams = make(map[string]proto.Message)

	return ret
}

func (i *DeployTxCall) Deploy(method string) error {

	msg := new(pb.DeployTx)
	msg.InitParams = make(map[string][]byte)

	for method, v := range i.InitParams {
		payload, err := proto.Marshal(v)
		if err != nil {
			return err
		}

		msg.InitParams[method] = payload
	}

	err := i.txcall(method, msg)
	if err != nil {
		return err
	}

	_, err = i.postHandling(method, call_deploy)
	return err
}
