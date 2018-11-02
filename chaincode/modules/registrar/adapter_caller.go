package registrar

import (
	"errors"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	pb "hyperledger.abchain.org/chaincode/modules/registrar/protos"
	"hyperledger.abchain.org/crypto"
)

type GeneralCall struct {
	*txgen.TxGenerator
}

const (
	Method_Registrar      = "REGISTRAR.PUBLICKEY"
	Method_AdminRegistrar = "REGISTRAR.PUBLICKEYDIRECT"
	Method_Revoke         = "REGISTRAR.REVOKEPK"
	Method_Active         = "REGISTRAR.ENABLEPK"
)

func (i *GeneralCall) AdminRegistrar(pk *crypto.PublicKey) error {

	msg := &pb.RegPublicKey{
		pk.PBMessage(),
		"",
	}

	_, err := i.Invoke(Method_AdminRegistrar, msg)
	return err
}

func (i *GeneralCall) Registrar(pk *crypto.PublicKey, region string) ([]byte, error) {

	msg := &pb.RegPublicKey{
		pk.PBMessage(),
		"",
	}

	_, err := i.Invoke(Method_Registrar, msg)

	if err != nil {
		return nil, err
	}

	return pk.RootFingerPrint, nil
}

func (i *GeneralCall) ActivePk(key []byte) error {

	msg := &pb.ActivePublicKey{
		key,
	}

	_, err := i.Invoke(Method_Active, msg)
	return err
}

func (i *GeneralCall) RevokePk(pk *crypto.PublicKey) error {

	msg := &pb.RevokePublicKey{
		pk.PBMessage(),
	}

	_, err := i.Invoke(Method_Revoke, msg)
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
	err = rpc.DecodeRPCResult(ret, d)
	if err != nil {
		return err, nil
	}

	return nil, ret
}

//don't need to call global methods in rpc
func (i *GeneralCall) Global() (error, *pb.RegGlobalData) {
	return errors.New("No implement"), nil
}
