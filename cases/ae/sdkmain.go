package main

import (
	"hyperledger.abchain.org/cases/ae/service"
	_ "hyperledger.abchain.org/chaincode/impl/hyfabric"
	_ "hyperledger.abchain.org/client/hyfabric/peerex"
)

func main() {

	// Start SDK Service
	service.StartService()
}
