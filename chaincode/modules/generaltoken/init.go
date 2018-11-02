package generaltoken

import (
	"errors"
	"github.com/golang/protobuf/proto"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
	"math/big"
)

const (
	deployName   = ":deploy"
	deployMethod = "GENERALTOKEN"
)

var DeployMethod = deployMethod

type DeployCall struct {
	*txgen.DeployTxCall
}

func (i *DeployCall) getParam() *pb.BaseToken {

	if v, ok := i.InitParams[deployMethod]; ok {
		return v.(*pb.BaseToken)
	} else {
		ret := new(pb.BaseToken)
		i.InitParams[deployMethod] = ret
		return ret
	}
}

func (i *DeployCall) Init(amount *big.Int) {

	msg := i.getParam()
	msg.TotalTokens = amount.Bytes()
}

func (i *DeployCall) Assign(to []byte, amount *big.Int) {

	msg := i.getParam()
	msg.Assigns = append(msg.Assigns, &pb.BaseToken_Assign{
		Recv:   txutil.NewAddressFromHash(to).PBMessage(),
		Amount: amount.Bytes(),
	})
}

type initHandler struct {
	msg pb.BaseToken
	TokenConfig
}

func CCDeployHandler(cfg TokenConfig) *initHandler {
	return &initHandler{TokenConfig: cfg}
}

func (h *initHandler) Msg() proto.Message { return &h.msg }

func (h *initHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := &h.msg

	token := h.NewTx(stub, parser.GetNounce())

	if err := token.Init(toAmount(msg.TotalTokens)); err != nil {
		return nil, err
	}

	for _, assign := range msg.Assigns {
		addr, err := txutil.NewAddressFromPBMessage(assign.Recv)
		if err != nil {
			return nil, err
		}

		_, err = token.Assign(addr.Hash, toAmount(assign.Amount))
		if err != nil {
			return nil, err
		}
	}

	return []byte("Ok"), nil

}

func (db *baseTokenTx) Init(total *big.Int) error {
	deploy := &pb.TokenGlobalData{}
	err := db.Get(deployName, deploy)

	if deploy.TotalTokens != nil {
		return errors.New("Can not re-deploy existed data")
	}

	err = db.Set(deployName, &pb.TokenGlobalData{total.Bytes(), total.Bytes()})
	if err != nil {
		return err
	}

	return nil
}
