package main

import (
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
	"hyperledger.abchain.org/chaincode"
)

func main() {
	cc := new(chaincode.AECC)
	cc.DebugMode = true
	<-ccutil.ExecuteCC(cc)
}
