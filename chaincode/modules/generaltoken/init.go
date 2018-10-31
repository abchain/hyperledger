package generaltoken

import (
	"errors"

	"hyperledger.abchain.org/chaincode/lib/caller"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	txutil "hyperledger.abchain.org/tx"
	"math/big"
)

const (
	deployName   = ":deploy"
	deployMethod = "GENERALTOKEN"
)

var DeployMethod = deployMethod

func CCDeploy(amount *big.Int, args []string) ([]string, error) {
	msg := &pb.BaseToken{
		amount.Bytes(),
		nil,
	}

	return rpc.BuildDeployArg(deployMethod, msg, args)
}

type CCDeployHandler string

func (tag CCDeployHandler) Call(stub interface{}, deployarg []byte) error {
	msg := &pb.BaseToken{}
	err := rpc.DecodeRPCResult(msg, deployarg)
	if err != nil {
		return err
	}

	cfg := StandardTokenConfig{}
	cfg.Tag = string(tag)
	cfg.Readonly = false
	token := cfg.NewTx(stub, []byte("DEPLOYTXS"))

	err = token.Init(toAmount(msg.TotalTokens))
	if err != nil {
		return err
	}

	for _, assign := range msg.Assigns {
		addr, err := txutil.NewAddressFromPBMessage(assign.Recv)
		if err != nil {
			return err
		}

		_, err = token.Assign(addr.Hash, toAmount(assign.Amount))
		if err != nil {
			return err
		}
	}

	return nil

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
