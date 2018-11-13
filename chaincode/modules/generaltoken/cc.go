package generaltoken

import (
	"encoding/base64"
	"hyperledger.abchain.org/chaincode/lib/runtime"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/shim"
	"math/big"
)

type TokenTx interface {
	nonce.TokenNonceTx
	Init(amount *big.Int) error
	Transfer(from []byte, to []byte, amount *big.Int) ([]byte, error)
	Assign(to []byte, amount *big.Int) ([]byte, error)
	Account(addr []byte) (error, *pb.AccountData_s)
	Global() (error, *pb.TokenGlobalData_s)
	//this is only used for inner call to register their address, have
	//no effect on the status of module
	TouchAddr([]byte) error
}

type TokenConfig interface {
	NewTx(shim.ChaincodeStubInterface, []byte) TokenTx
}

//the local config must provide both a executable interface and the sub-config (corresponding its sub interface)
//for local handler building
type LocalConfig interface {
	TokenConfig
	Nonce() nonce.NonceConfig
}

type StandardTokenConfig struct {
	Root string
	*runtime.Config
	NonceCfg nonce.NonceConfig
}

func NewConfig(tag string) *StandardTokenConfig {
	cfg := runtime.NewConfig()
	nccfg := nonce.NewConfig(tag)
	nccfg.Config = cfg
	return &StandardTokenConfig{tx_tag_prefix + tag, cfg, nccfg}
}

type baseTokenTx struct {
	*runtime.ChaincodeRuntime
	nonce      []byte
	tokenNonce nonce.TokenNonceTx
}

const (
	tx_tag_prefix = "GenToken_"
)

func (cfg *StandardTokenConfig) NewTx(stub shim.ChaincodeStubInterface, nc []byte) TokenTx {

	return &baseTokenTx{runtime.NewRuntime(cfg.Root, stub, cfg.Config), nc, cfg.NonceCfg.NewTx(stub, nc)}

}

func (cfg *StandardTokenConfig) Nonce() nonce.NonceConfig { return cfg.NonceCfg }

type InnerInvokeConfig struct {
	txgen.InnerChaincode
}

func (c InnerInvokeConfig) NewTx(stub shim.ChaincodeStubInterface, nc []byte) TokenTx {
	return NewFullGeneralCall(c.NewInnerTxInterface(stub, nc))
}

func toAmount(a []byte) *big.Int {

	if a == nil {
		return big.NewInt(0)
	}

	return big.NewInt(0).SetBytes(a)
}

func addrToKey(h []byte) string {
	return base64.RawURLEncoding.EncodeToString(h)
}
