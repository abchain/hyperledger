package generaltoken

import (
	"encoding/base64"
	"hyperledger.abchain.org/chaincode/lib/runtime"
	"hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/shim"
	"math/big"
)

type nonceResolver interface {
	Nonce(stub shim.ChaincodeStubInterface, nonce []byte) nonce.TokenNonceTx
}

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

type StandardTokenConfig struct {
	Tag      string
	Readonly bool
	Nonce    nonceResolver
}

type InnerNonceResolver struct {
	*nonce.StandardNonceConfig
}

func (i InnerNonceResolver) Nonce(stub shim.ChaincodeStubInterface, _ []byte) nonce.TokenNonceTx {
	return i.NewTx(stub)
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
	rootname := tx_tag_prefix + cfg.Tag
	nresolver := cfg.Nonce
	if nresolver == nil {
		nresolver = InnerNonceResolver{&nonce.StandardNonceConfig{cfg.Tag, cfg.Readonly}}
	}

	return &baseTokenTx{runtime.NewRuntime(rootname, stub, cfg.Readonly), nc, nresolver.Nonce(stub, nc)}
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
