package registrar

import (
	"errors"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	pb "hyperledger.abchain.org/chaincode/modules/registrar/protos"
	"hyperledger.abchain.org/core/crypto"
)

type GeneralCall struct {
	*txgen.TxGenerator
}

const (
	Method_Registrar      = "REGISTRAR.PUBLICKEY"
	Method_AdminRegistrar = "REGISTRAR.PUBLICKEYDIRECT"
	Method_Revoke         = "REGISTRAR.REVOKEPK"
	Method_Active         = "REGISTRAR.ENABLEPK"
	Method_Init           = "REGISTRAR.INIT"
)

func (i *GeneralCall) AdminRegistrar(pk *crypto.PublicKey) error {

	msg := &pb.RegPublicKey{
		pk.PBMessage(),
		"",
	}

	err := i.Invoke(Method_AdminRegistrar, msg)
	return err
}

func (i *GeneralCall) Registrar(pk *crypto.PublicKey, region string) ([]byte, error) {

	msg := &pb.RegPublicKey{
		pk.PBMessage(),
		"",
	}

	err := i.Invoke(Method_Registrar, msg)

	if err != nil {
		return nil, err
	}

	return pk.RootFingerPrint, nil
}

func (i *GeneralCall) ActivePk(key []byte) error {

	msg := &pb.ActivePublicKey{
		key,
	}

	err := i.Invoke(Method_Active, msg)
	return err
}

func (i *GeneralCall) RevokePk(pk *crypto.PublicKey) error {

	msg := &pb.RevokePublicKey{
		pk.PBMessage(),
	}

	err := i.Invoke(Method_Revoke, msg)
	return err
}

func (i *GeneralCall) Pubkey(key []byte) (error, *pb.RegData) {
	msg := &pb.ActivePublicKey{
		key,
	}

	d, err := i.Query(Method_Registrar, msg) //re-use the method name for invoke
	if err != nil {
		return err, nil
	}

	ret := &pb.RegData{}
	err = txgen.SyncQueryResult(ret, d)
	if err != nil {
		return err, nil
	}

	return nil, ret
}

//don't need to call global methods in rpc
func (i *GeneralCall) Global() (error, *pb.RegGlobalData) {
	return errors.New("No implement"), nil
}
