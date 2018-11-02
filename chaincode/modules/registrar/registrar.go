package registrar

import (
	pb "hyperledger.abchain.org/chaincode/modules/registrar/protos"
	"hyperledger.abchain.org/crypto"
	"hyperledger.abchain.org/utils"

	"errors"
)

func (db *registrarTx) registrar(pk *crypto.PublicKey, region string, enable bool) error {

	if !pk.Index.IsUint64() || pk.Index.Uint64() != 0 {
		return errors.New("Can only register root publickey")
	}

	regkey := registrarKey(pk)
	data := &pb.RegData{}
	err := db.Get(regkey, data)

	if err != nil {
		return err
	}

	if data.Pk != nil {
		return errors.New("Public key has been reg")
	}

	t, _ := db.stub.GetTxTime()
	data = &pb.RegData{
		pk.PBMessage(),
		db.stub.GetTxID(),
		utils.CreatePBTimestamp(t), region, enable, nil,
	}

	return db.Set(regkey, data)
}

func (db *registrarTx) AdminRegistrar(pk *crypto.PublicKey) error {

	attr, err := db.stub.GetCallerAttribute(db.RegionAttrName)
	if err != nil {
		return err
	}

	return db.registrar(pk, string(attr), true)
}

func (r *registrarTx) Registrar(pk *crypto.PublicKey, region string) ([]byte, error) {
	err := r.registrar(pk, region, false)
	if err != nil {
		return nil, err
	}

	return pk.RootFingerPrint, nil
}

func (db *registrarTx) ActivePk(key []byte) error {

	regkey := registrarQueryKey(key)
	data := &pb.RegData{}
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

func (db *registrarTx) RevokePk(pk *crypto.PublicKey) error {
	return errors.New("No implement")
}

func (db *registrarTx) Pubkey(key []byte) (error, *pb.RegData) {

	regkey := registrarQueryKey(key)
	data := &pb.RegData{}
	err := db.Get(regkey, data)

	if err != nil {
		return err, nil
	}

	if data.Pk == nil {
		return errors.New("The public key is not reg yet"), nil
	}

	return nil, data
}

func (db *registrarTx) Global() (error, *pb.RegGlobalData) {

	data := &pb.RegGlobalData{}
	err := db.Get(deployName, data)

	if err != nil {
		return err, nil
	}

	if data.DeployFlag == nil {
		return errors.New("Not deploy yet"), nil
	}

	return nil, data
}
