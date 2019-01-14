package multitoken

import (
	"hyperledger.abchain.org/chaincode/modules/generaltoken"
	"hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/shim"
)

type MultiTokenTx interface {
	GetToken(string) generaltoken.TokenTx
	Create(string) error
}
