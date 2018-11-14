package service

import (
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
	"hyperledger.abchain.org/chaincode/lib/caller"
)

var ccCaller *rpc.ChaincodeAdapter

func init() {
	ccCaller = rpc.NewLocalChaincode(chaincode.NewChaincode(true))
}
