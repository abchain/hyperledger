package main

import (
	"fmt"
	log "github.com/abchain/fabric/peerex/logging"
	"hyperledger.abchain.org/cases/ae/service"
	"hyperledger.abchain.org/config"
)

var Version = "v0.1.0"

var logger = log.InitLogger("service")

func main() {

	// Read config
	globalConfig := &config.GlobalConfig{ConfigFileName: "ae_conf"}
	err := globalConfig.InitGlobal()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set log level, these will be applied when StarService is called
	log.SetLogLevel("")
	log.SetLogFormat(`%{color}%{time:15:04:05.000} %{level:.4s} [%{module:.6s}] %{shortfile} %{shortfunc} ▶ %{message}%{color:reset}`)

	// Print Version
	logger.Infof("Atomic Power SDK: %s", Version)

	// Start SDK Service
	service.StartService()

}
