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
	parser := txhandle.GenerateTxArgParser(aecc.CollectiveTxs.Map())
	parser[chaincode.CC_BATCH+"@"+chaincode.CC_NAME] = txhandle.BatchArgParser(chaincode.CC_NAME, parser)
	parser[client.TxErrorEventName] = client.TxErrorParser("Local invoke error:")
	blockchain.SetParsers(parser)
}
