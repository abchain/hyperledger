package multitoken

import (
	"hyperledger.abchain.org/chaincode/lib/runtime"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"hyperledger.abchain.org/chaincode/modules/generaltoken"
	"hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
	"hyperledger.abchain.org/chaincode/shim"
)

//Currying: except for createToken, most of the tx in multitoken formed by two continuous calling:
//GetToken and then one of the methods in the returned TokenTx
type MultiTokenTx interface {
	GetToken(string) (generaltoken.TokenTx, error)
}

type baseMultiTokenTx struct {
	*runtime.ChaincodeRuntime
	nonce      []byte
	tokenNonce nonce.TokenNonceTx
}

type TokenConfig interface {
	NewTx(shim.ChaincodeStubInterface, []byte) MultiTokenTx
}

type StandardTokenConfig struct {
	Root string
	*runtime.Config
	NonceCfg nonce.NonceConfig
}

func (cfg *StandardTokenConfig) NewTx(stub shim.ChaincodeStubInterface, nc []byte) MultiTokenTx {

	return &baseMultiTokenTx{runtime.NewRuntime(cfg.Root, stub, cfg.Config), nc, cfg.NonceCfg.NewTx(stub, nc)}

}

const (
	tx_tag_prefix = "GenMultiToken_"
)

func NewConfig(tag string) *StandardTokenConfig {
	cfg := runtime.NewConfig()
	nccfg := nonce.NewConfig(tag)
	nccfg.Config = cfg
	return &StandardTokenConfig{tx_tag_prefix + tag, cfg, nccfg}
}

func ConfigFromToken(tcf *generaltoken.StandardTokenConfig) *StandardTokenConfig {
	return &StandardTokenConfig{tcf.Root, tcf.Config, tcf.NonceCfg}
}

type InnerInvokeConfig struct {
	txgen.InnerChaincode
}

func (c InnerInvokeConfig) NewTx(stub shim.ChaincodeStubInterface, nc []byte) MultiTokenTx {
	return &GeneralCall{c.NewInnerTxInterface(stub, nc)}
}
