package main

import (
	log "github.com/abchain/fabric/peerex/logging"
	"hyperledger.abchain.org/cases/ae/service"
	"os"
)

func main() {

	log.SetLogFormat(`%{color}%{time:15:04:05.000} %{level:.4s} [%{module:.6s}] %{shortfile} %{shortfunc} ▶ %{message}%{color:reset}`)
	log.SetBackend(os.Stderr, "", 0)

	// Start SDK Service
	service.StartService()
}
