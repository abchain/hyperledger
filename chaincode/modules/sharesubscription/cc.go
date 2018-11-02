package subscription

import (
	"hyperledger.abchain.org/chaincode/lib/state"
	token "hyperledger.abchain.org/chaincode/modules/generaltoken"
	pb "hyperledger.abchain.org/chaincode/modules/sharesubscription/protos"
	"hyperledger.abchain.org/chaincode/shim"
	"hyperledger.abchain.org/crypto"
	"math/big"
)

type ContractTx interface {
	New(map[string]uint32, *crypto.PublicKey) ([]byte, error)                               //return contract address
	Redeem(conaddr []byte, addr []byte, amount *big.Int, redeemAddr []byte) ([]byte, error) //return noncekey in token
	Query(addr []byte) (error, *pb.Contract)
	QueryOne(conaddr []byte, addr []byte) (error, *pb.Contract)
}

type ContractConfig interface {
	NewTx(shim.ChaincodeStubInterface, []byte) ContractTx
}

type StandardContractConfig struct {
	Tag      string
	Readonly bool
	TokenCfg token.TokenConfig
}

type baseContractTx struct {
	state.StateMap
	nonce []byte
	stub  shim.ChaincodeStubInterface
	token token.TokenTx
}

const (
	contract_tag_prefix = "Subscription_"
)

func (cfg *StandardContractConfig) NewTx(stub shim.ChaincodeStubInterface, nonce []byte) ContractTx {
	rootname := contract_tag_prefix + cfg.Tag

	return &baseContractTx{state.NewShimMap(rootname, stub, cfg.Readonly), nonce,
		stub, cfg.TokenCfg.NewTx(stub, nonce)}
}
