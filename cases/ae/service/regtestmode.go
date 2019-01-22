package service

import (
	"hyperledger.abchain.org/applications/blockchain"
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
	txhandle "hyperledger.abchain.org/chaincode/lib/txhandle"
	"hyperledger.abchain.org/client/local"
)

func init() {
	aecc := chaincode.NewChaincode(true)
	client.AddChaincode("local", aecc)

	//also build txparser from chaincode ...
	blockchain.SetParsers(txhandle.GenerateTxArgParser(aecc.CollectiveTxs.Map()))
}
