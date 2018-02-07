package generaltoken

import (
	"encoding/base64"
	"github.com/abchain/fabric/core/chaincode/shim"
	"hyperledger.abchain.org/chaincode/generaltoken/nonce"
	pb "hyperledger.abchain.org/chaincode/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/lib/state"
	"math/big"
)

type TokenTx interface {
	nonce.TokenNonceTx
	Init(amount *big.Int) error
	Transfer(from []byte, to []byte, amount *big.Int) ([]byte, error)
	Assign(to []byte, amount *big.Int) ([]byte, error)
	Account(addr []byte) (error, *pb.AccountData)
	Global() (error, *pb.TokenGlobalData)
}

type TokenConfig interface {
	NewTx(shim.ChaincodeStubInterface, []byte) TokenTx
}

//integrate both nonce and token
type StandardTokenConfig struct {
	nonce.StandardNonceConfig
}

type baseTokenTx struct {
	state.StateMap
	nonce      []byte
	stub       shim.ChaincodeStubInterface
	tokenNonce nonce.TokenNonceTx
}

const (
	tx_tag_prefix = "GenToken_"
)

func (cfg *StandardTokenConfig) NewTx(stub shim.ChaincodeStubInterface, nonce []byte) TokenTx {
	rootname := tx_tag_prefix + cfg.Tag

	return &baseTokenTx{state.NewShimMap(rootname, stub, cfg.Readonly), nonce, stub,
		cfg.StandardNonceConfig.NewTx(stub)}
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
