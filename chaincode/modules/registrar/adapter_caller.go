package registrar

import (
	"errors"

	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	pb "hyperledger.abchain.org/chaincode/modules/registrar/protos"
	"hyperledger.abchain.org/core/crypto"
	protos "hyperledger.abchain.org/protos"
)

type GeneralCall struct {
	txgen.TxCaller
}

const (
	Method_Registrar      = "REGISTRAR.PUBLICKEY"
	Method_AdminRegistrar = "REGISTRAR.PUBLICKEYDIRECT"
	Method_Revoke         = "REGISTRAR.REVOKEPK"
	Method_Active         = "REGISTRAR.ENABLEPK"
	Method_Init           = "REGISTRAR.INIT"
)

func (i *GeneralCall) AdminRegistrar(pkbytes []byte) error {

	msg := &pb.RegPublicKey{
		PkBytes: pkbytes,
		Region:  "",
	}

	err := i.Invoke(Method_AdminRegistrar, msg)
	return err
}

func (i *GeneralCall) Registrar(pkbytes []byte, region string) error {

	msg := &pb.RegPublicKey{
		PkBytes: pkbytes,
		Region:  "",
	}

	err := i.Invoke(Method_Registrar, msg)
	return err
}

func (i *GeneralCall) ActivePk(key []byte) error {

	msg := &pb.ActivePublicKey{
		key,
	}

	err := i.Invoke(Method_Active, msg)
	return err
}

func (i *GeneralCall) RevokePk(pk crypto.Verifier) error {

	msg := &pb.RevokePublicKey{
		pk.PBMessage().(*protos.PublicKey),
	}

	err := i.Invoke(Method_Revoke, msg)
	return err
}

func (i *GeneralCall) InitDebugMode() error {

	return i.Invoke(Method_Init, &pb.Settings{true, "", ""})
}

func (i *GeneralCall) Init(enablePrivilege bool, managePriv string, regPriv string) error {

	return i.Invoke(Method_Init, &pb.Settings{enablePrivilege, managePriv, regPriv})
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
