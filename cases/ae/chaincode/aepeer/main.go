package main

import (
	node "github.com/abchain/fabric/node/start"
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
	"hyperledger.abchain.org/chaincode/impl/yafabric"

	"github.com/abchain/fabric/core/embedded_chaincode/api"
	"github.com/abchain/fabric/node/legacy"
)

func main() {

	adapter := &legacynode.LegacyEngineAdapter{}

	reg := func() error {
		cc := fabric_impl.GenYAfabricCC(chaincode.NewChaincode(true))
		if err := api.RegisterECC(&api.EmbeddedChaincode{"aecc", cc}); err != nil {
			return err
		}

		if err := adapter.Init(); err != nil {
			return err
		}

		return nil
	}

	node.RunNode(&node.NodeConfig{PostRun: reg, Schemes: adapter.Scheme})

}
