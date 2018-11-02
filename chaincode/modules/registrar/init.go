package registrar

import (
	"errors"
	"github.com/golang/protobuf/proto"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	pb "hyperledger.abchain.org/chaincode/modules/registrar/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
)

const (
	deployName    = ":deploy"
	DeployMethod  = "REGISTRAR"
	debugModeFlag = 'D'
)

type DeployCall struct {
	*txgen.DeployTxCall
}

func (i *DeployCall) Init(enablePrivilege bool, managePriv string, regPriv string) {

	i.InitParams[DeployMethod] = &pb.Settings{enablePrivilege, managePriv, regPriv}
}

func (i *DeployCall) InitDebugMode() {

	i.InitParams[DeployMethod] = &pb.Settings{true, "", ""}
}

type initHandler struct {
	msg pb.Settings
	RegistrarConfig
}

func CCDeployHandler(cfg RegistrarConfig) *initHandler {
	return &initHandler{RegistrarConfig: cfg}
}

func (h *initHandler) Msg() proto.Message { return &h.msg }

func (h *initHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := &h.msg

	reg := h.NewTx(stub)

	err := reg.Init(!msg.DebugMode, msg.AdminPrivilege, msg.RegPrivilege)
	if err != nil {
		return nil, err
	}

	return []byte("Ok"), nil
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
