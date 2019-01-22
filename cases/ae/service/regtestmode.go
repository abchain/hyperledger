package service

import (
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
	"hyperledger.abchain.org/client/local"
)

func init() {
	client.AddChaincode("local", chaincode.NewChaincode(true))
}
