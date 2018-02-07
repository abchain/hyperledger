package service

import (
	"github.com/abchain/fabric/core/chaincode/shim"
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
	"hyperledger.abchain.org/chaincode/lib/caller"
)

var ccCaller *rpc.ChaincodeAdapter

func init() {
	ccCaller = &rpc.ChaincodeAdapter{shim.NewMockStub("RegTestMode", &chaincode.AECC{DebugMode: true}), nil}
}
