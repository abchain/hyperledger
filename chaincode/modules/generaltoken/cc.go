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
	Init(amount *big.Int) error
	Transfer(from []byte, to []byte, amount *big.Int) ([]byte, error)
	Assign(to []byte, amount *big.Int) ([]byte, error)
	Account(addr []byte) (error, *pb.AccountData_s)
	Global() (error, *pb.TokenGlobalData_s)
	//this is only used for inner call to register their address, have
	//no effect on the status of module
	TouchAddr([]byte) error
}

type TokenTxExt interface {
	TokenTx
	nonce.TokenNonceTx
}

type TokenConfig interface {
	NewTx(shim.ChaincodeStubInterface, []byte) TokenTx
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

func NewTokenTxImpl(rt *runtime.ChaincodeRuntime, nc []byte, ncTx nonce.TokenNonceTx) *baseTokenTx {
	return &baseTokenTx{rt, nc, ncTx}
}

const (
	tx_tag_prefix = "GenToken_"
)

func (cfg *StandardTokenConfig) NewTx(stub shim.ChaincodeStubInterface, nc []byte) TokenTx {

	return NewTokenTxImpl(runtime.NewRuntime(cfg.Root, stub, cfg.Config), nc, cfg.NonceCfg.NewTx(stub, nc))

}

//func (cfg *StandardTokenConfig) Nonce() nonce.NonceConfig { return cfg.NonceCfg }

type InnerInvokeConfig struct {
	txgen.InnerChaincode
}

func (c InnerInvokeConfig) NewTx(stub shim.ChaincodeStubInterface, nc []byte) TokenTx {
	return &GeneralCall{c.NewInnerTxInterface(stub, nc)}
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
