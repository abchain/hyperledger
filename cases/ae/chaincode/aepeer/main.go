package main

import (
	"github.com/abchain/fabric/peerex/node"
	"hyperledger.abchain.org/cases/ae/chaincode/cc"

	"github.com/abchain/fabric/core/embedded_chaincode/api"
)

func main() {

	reg := func() error {
		return api.RegisterECC(&api.EmbeddedChaincode{"aecc", new(chaincode.AECC)})
	}

	node.RunNode(&node.NodeConfig{PostRun: reg})

}
