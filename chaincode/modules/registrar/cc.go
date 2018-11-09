package registrar

import (
	"encoding/base64"
	"hyperledger.abchain.org/chaincode/lib/runtime"
	pb "hyperledger.abchain.org/chaincode/modules/registrar/protos"
	"hyperledger.abchain.org/chaincode/shim"
	"hyperledger.abchain.org/core/crypto"
)

type RegistrarTx interface {
	Init(enablePrivilege bool, managePriv string, regPriv string) error
	AdminRegistrar(pk *crypto.PublicKey) error
	Registrar(pk *crypto.PublicKey, region string) ([]byte, error)
	ActivePk(key []byte) error
	RevokePk(pk *crypto.PublicKey) error
	Pubkey(key []byte) (error, *pb.RegData)
	Global() (error, *pb.RegGlobalData)
}

type RegistrarConfig interface {
	NewTx(shim.ChaincodeStubInterface) RegistrarTx
}

type StandardRegistrarConfig struct {
	Tag                string
	Readonly           bool
	ManagePrivAttrName string
	RegionAttrName     string
}

const (
	reg_tag_prefix = "Registrar_"
)

type registrarTx struct {
	runtime.StateMap_Legacy
	stub shim.ChaincodeStubInterface
	*StandardRegistrarConfig
}

func (cfg *StandardRegistrarConfig) NewTx(stub shim.ChaincodeStubInterface) RegistrarTx {
	rootname := reg_tag_prefix + cfg.Tag

	return &registrarTx{runtime.NewShimMapLegacy(rootname, stub, cfg.Readonly), stub, cfg}
}

func registrarQueryKey(key []byte) string {
	return base64.RawURLEncoding.EncodeToString(key)
}

func registrarKey(pk *crypto.PublicKey) string {
	return registrarQueryKey(pk.RootFingerPrint)
}
