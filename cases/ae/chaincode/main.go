package main

import (
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
	"hyperledger.abchain.org/chaincode"
)

func main() {
	<-ccutil.ExecuteCC(new(chaincode.AECC))
}
