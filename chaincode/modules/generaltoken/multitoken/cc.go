package multitoken

import (
	"hyperledger.abchain.org/chaincode/lib/runtime"
	"hyperledger.abchain.org/chaincode/modules/generaltoken"
	"hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
	_ "hyperledger.abchain.org/chaincode/shim"
	"math/big"
)

//Currying: except for createToken, most of the tx in multitoken formed by two continuous calling:
//GetToken and then one of the methods in the returned TokenTx
type MultiTokenTx interface {
	GetToken(string) (generaltoken.TokenTxCore, error)
	CreateToken(string, *big.Int) error
}

type baseMultiTokenTx struct {
	*runtime.ChaincodeRuntime
	nonce      []byte
	tokenNonce nonce.TokenNonceTx
}

type TokenConfig interface {
	NewTx(shim.ChaincodeStubInterface, []byte) MultiTokenTx
}
