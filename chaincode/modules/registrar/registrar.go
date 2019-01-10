package registrar

import (
	"errors"

	"hyperledger.abchain.org/chaincode/impl"
	pb "hyperledger.abchain.org/chaincode/modules/registrar/protos"
	"hyperledger.abchain.org/core/crypto"
)

const (
	deployName    = ":deploy"
	debugModeFlag = 'D'
)

func (db *registrarTx) registrar(pkbyte []byte, region string, enable bool) error {

	pk, err := crypto.PublicKeyFromBytes(pkbyte)
	if err != nil {
		return err
	}

	if len(pk.GetRootFingerPrint()) != 0 {
		return errors.New("Can only register root publickey")
	}

	regkey := registrarKey(pk)
	data := &pb.RegData_s{}
	err = db.Get(regkey, data)

	if err != nil {
		return err
	}

	data.PkBytes = pkbyte
	data.RegTxid = db.stub.GetTxID()
	data.RegTs, _ = db.stub.GetTxTime()
	data.Enabled = enable
	data.Region = region

	if data.Pk != nil {
		return errors.New("Public key has been reg")
	}

	return db.Set(regkey, data)
}

func (db *registrarTx) AdminRegistrar(pkbyte []byte) error {

	attrif, err := impl.GetCallerAttributes(db.stub)
	if err != nil {
		return err
	}

	attr, err := attrif.GetCallerAttribute(db.RegionAttrName)
	if err != nil {
		return err
	}

	return db.registrar(pkbyte, string(attr), true)
}

func (r *registrarTx) Registrar(pkbyte []byte, region string) error {
	return r.registrar(pkbyte, region, false)

}

func (db *registrarTx) ActivePk(key []byte) error {

	regkey := registrarQueryKey(key)
	data := &pb.RegData_s{}
	err := db.Get(regkey, data)

	if err != nil {
		return err
	}

	if data.Pk == nil {
		return errors.New("Public key has not been reg")
	}

	if data.Enabled {
		return nil
	} else {
		data.Enabled = true
		return db.Set(regkey, data)
	}

}

func (db *registrarTx) RevokePk(pk crypto.Verifier) error {
	return errors.New("No implement")
}

func (db *registrarTx) Init(enablePrivilege bool, managePriv string, regPriv string) error {

	deploy := &pb.RegGlobalData_s{}
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

func (db *registrarTx) pubkey(key []byte) (error, *pb.RegData_s) {
	regkey := registrarQueryKey(key)
	data := &pb.RegData_s{}
	err := db.Get(regkey, data)

	if err != nil {
		return err, nil
	}

	if data.Pk == nil {
		return errors.New("The public key is not reg yet"), nil
	}

	return nil, data
}

func (db *registrarTx) Pubkey(key []byte) (error, *pb.RegData) {

	err, data := db.pubkey(key)
	if err != nil {
		return err, nil
	}

	return nil, data.ToPB()
}

func (db *registrarTx) Global() (error, *pb.RegGlobalData) {

	data := &pb.RegGlobalData_s{}
	err := db.Get(deployName, data)

	if err != nil {
		return err, nil
	}

	if data.DeployFlag == nil {
		return errors.New("Not deploy yet"), nil
	}

	return nil, data.ToPB()
}
