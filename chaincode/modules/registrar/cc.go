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
	AdminRegistrar(pkbyte []byte) error
	Registrar(pkbyte []byte, region string) error
	ActivePk(key []byte) error
	RevokePk(pk crypto.Verifier) error
	Pubkey(key []byte) (error, *pb.RegData)
	Global() (error, *pb.RegGlobalData)
}

type RegistrarTxExt interface {
	RegistrarTx
	pubkey(key []byte) (error, *pb.RegData_s)
}

type RegistrarConfig interface {
	NewTx(shim.ChaincodeStubInterface) RegistrarTxExt
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
	runtime.StateMap
	stub shim.ChaincodeStubInterface
	*StandardRegistrarConfig
}

func (cfg *StandardRegistrarConfig) NewTx(stub shim.ChaincodeStubInterface) RegistrarTxExt {
	rootname := reg_tag_prefix + cfg.Tag

	return &registrarTx{runtime.NewShimMap(rootname, stub, cfg.Readonly), stub, cfg}
}

func registrarQueryKey(key []byte) string {
	if len(key) > crypto.PUBLICKEY_FINGERPRINT_LEN {
		key = key[:crypto.PUBLICKEY_FINGERPRINT_LEN]
	}
	return base64.RawURLEncoding.EncodeToString(key)
}

func registrarKey(pk crypto.Verifier) string {

	return registrarQueryKey(pk.Digest())
}
