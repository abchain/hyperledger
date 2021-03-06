package subscription

import (
	"hyperledger.abchain.org/chaincode/lib/runtime"
	"hyperledger.abchain.org/chaincode/modules/addrspace"
	token "hyperledger.abchain.org/chaincode/modules/generaltoken"
	pb "hyperledger.abchain.org/chaincode/modules/sharesubscription/protos"
	"hyperledger.abchain.org/chaincode/shim"
	"math/big"
)

type ContractTx interface {
	New_C(addrs [][]byte, ratios []int) ([]byte, error) //return contract address
	Redeem_C(conaddr []byte, amount *big.Int, redeemAddrs [][]byte) (*pb.RedeemResponse, error)
	Query_C(addr []byte) (error, *pb.Contract_s)
	QueryOne_C(conaddr, addr []byte) (error, *pb.Contract_s)
}

type ContractConfig interface {
	NewTx(shim.ChaincodeStubInterface, []byte) ContractTx
}

type StandardContractConfig struct {
	Root string
	*runtime.Config
	TokenCfg token.TokenConfig
	AddrCfg  addrspace.AddrSpaceConfig
}

func NewConfig(tag string) *StandardContractConfig {
	cfg := runtime.NewConfig()
	tkcfg := token.NewConfig(tag)
	tkcfg.Config = cfg

	return &StandardContractConfig{contract_tag_prefix + tag,
		cfg, tkcfg, addrspace.DummyImplCfg()}
}

type baseContractTx struct {
	*runtime.ChaincodeRuntime
	nonce    []byte
	token    token.TokenTx
	addrutil addrspace.AddressSpace
}

const (
	contract_tag_prefix = "Subscription_"
	contract_auth_tag   = "Auth"
)

func (cfg *StandardContractConfig) NewTx(stub shim.ChaincodeStubInterface, nonce []byte) ContractTx {

	return &baseContractTx{runtime.NewRuntime(cfg.Root, stub, cfg.Config),
		nonce, cfg.TokenCfg.NewTx(stub, nonce), cfg.AddrCfg.NewTx(stub, nonce)}
}
