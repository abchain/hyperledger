package subscription

import (
	"hyperledger.abchain.org/chaincode/lib/runtime"
	token "hyperledger.abchain.org/chaincode/modules/generaltoken"
	pb "hyperledger.abchain.org/chaincode/modules/sharesubscription/protos"
	"hyperledger.abchain.org/chaincode/shim"
	"hyperledger.abchain.org/core/crypto"
	"math/big"
)

type ContractTx interface {
	New(map[string]int32, *crypto.PublicKey) ([]byte, error)                                //return contract address
	Redeem(conaddr []byte, addr []byte, amount *big.Int, redeemAddr []byte) ([]byte, error) //return noncekey in token
	Query(addr []byte) (error, *pb.Contract_s)
	QueryOne(conaddr []byte, addr []byte) (error, *pb.Contract_s)
}

type ContractConfig interface {
	NewTx(shim.ChaincodeStubInterface, []byte) ContractTx
}

type StandardContractConfig struct {
	Root string
	*runtime.Config
	TokenCfg token.TokenConfig
}

func NewConfig(tag string) *StandardContractConfig {
	cfg := runtime.NewConfig()
	tkcfg := token.NewConfig(tag)
	tkcfg.Config = cfg
	return &StandardContractConfig{contract_tag_prefix + tag, cfg, tkcfg}
}

type baseContractTx struct {
	*runtime.ChaincodeRuntime
	nonce []byte
	token token.TokenTx
}

const (
	contract_tag_prefix = "Subscription_"
)

func (cfg *StandardContractConfig) NewTx(stub shim.ChaincodeStubInterface, nonce []byte) ContractTx {
	return &baseContractTx{runtime.NewRuntime(cfg.Root, stub, cfg.Config), nonce, cfg.TokenCfg.NewTx(stub, nonce)}
}
