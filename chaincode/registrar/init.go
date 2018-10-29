package registrar

import (
	"errors"
	"hyperledger.abchain.org/chaincode/lib/caller"
	pb "hyperledger.abchain.org/chaincode/registrar/protos"
)

const (
	deployName    = ":deploy"
	DeployMethod  = "REGISTRAR"
	debugModeFlag = 'D'
)

func CCDeploy(managePriv string, regPriv string, args []string) ([]string, error) {
	msg := &pb.Settings{false, managePriv, regPriv}

	return rpc.BuildDeployArg(DeployMethod, msg, args)
}

func CCDeployDebugMode(args []string) ([]string, error) {
	msg := &pb.Settings{true, "", ""}

	return rpc.BuildDeployArg(DeployMethod, msg, args)
}

type CCDeployHandler string

func (tag CCDeployHandler) Call(stub interface{}, deployarg []byte) error {
	msg := &pb.Settings{}
	err := rpc.DecodeRPCResult(msg, deployarg)
	if err != nil {
		return err
	}

	cfg := StandardRegistrarConfig{}
	cfg.Tag = string(tag)
	cfg.Readonly = false
	reg := cfg.NewTx(stub)

	err = reg.Init(!msg.DebugMode, msg.AdminPrivilege, msg.RegPrivilege)
	if err != nil {
		return err
	}

	return nil

}

func (db *registrarTx) Init(enablePrivilege bool, managePriv string, regPriv string) error {

	deploy := &pb.RegGlobalData{}
	err := db.Get(deployName, deploy)

	if err != nil {
		return err
	}

	if deploy.DeployFlag != nil {
		return errors.New("Can not re-deploy existed data")
	}

	deploy.AdminPrivilege = managePriv
	deploy.RegPrivilege = regPriv

	if enablePrivilege {
		deploy.DeployFlag = []byte{'O', 'K'}
	} else {
		deploy.DeployFlag = []byte{debugModeFlag}
	}

	return db.Set(deployName, deploy)

}
