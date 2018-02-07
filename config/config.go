package config

import (
	log "github.com/abchain/fabric/peerex/logging"
	vieprWrapper "github.com/abchain/fabric/peerex/viper"
)

var defaultConfigFileName = "sdk_conf"

var logger = log.InitLogger("CONFIG")

var viper = vieprWrapper.New()
var Viper = viper

type GlobalConfig struct {
	ConfigFileName string
	ConfigPath     []string
}

func (g *GlobalConfig) InitGlobal() error {

	// Init ConfigPath
	if g.ConfigPath == nil {
		g.ConfigPath = make([]string, 1, 10)
		g.ConfigPath[0] = "."
	}

	// Init ConfigFileName
	if g.ConfigFileName == "" {
		g.ConfigFileName = defaultConfigFileName
	}

	// Set Default Value
	err := g.SetDefaultValue()
	if err != nil {
		return err
	}

	// Load config
	err = g.LoadConfig()
	if err != nil {
		return err
	}

	return nil
}

func (g *GlobalConfig) LoadConfig() error {
	for _, c := range g.ConfigPath {
		viper.AddConfigPath(c)
	}

	viper.SetConfigName(g.ConfigFileName) // Name of config file (without extension)

	logger.Debugf("ConfigPath: %v", g.ConfigPath)
	logger.Debugf("ConfigFileName: %v", g.ConfigFileName)

	return viper.ReadInConfig() // Find and read the config file
}

func (g *GlobalConfig) SetDefaultValue() error {

	// Debug
	logging := map[string]interface{}{}
	logging["level"] = "debug"
	viper.SetDefault("logging", logging)

	// Wallet
	wallet := map[string]interface{}{}
	wallet["path"] = "data"
	wallet["filename"] = "simplewallet.dat"
	viper.SetDefault("wallet", wallet)

	// Local RPC Server
	service := map[string]interface{}{}
	service["host"] = "localhost"
	service["port"] = "7080"
	viper.SetDefault("service", service)

	// gRPC Server
	grpc := map[string]interface{}{}
	grpc["server"] = "example.abchain.org:8000"
	grpc["username"] = "nobody"
	grpc["chaincode"] = "examplechain"
	grpc["tlsenabled"] = false
	grpc["certfile"] = "ca.crt"
	grpc["path"] = "data"
	viper.SetDefault("grpc", grpc)

	// REST Server
	rest := map[string]interface{}{}
	rest["server"] = "http://localhost:8080"
	viper.SetDefault("rest", rest)

	// Setting
	setting := map[string]interface{}{}
	setting["offline"] = false
	viper.SetDefault("setting", setting)

	return nil
}
