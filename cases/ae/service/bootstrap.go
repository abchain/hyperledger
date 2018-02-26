package service

import (
	"hyperledger.abchain.org/asset/wallet"
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
	"hyperledger.abchain.org/client"
	"hyperledger.abchain.org/config"
	"path/filepath"
)

const (
	defaultCCDeployName = "aecc"
)

var (
	defaultWallet    wallet.Wallet
	defaultRpcConfig *client.RpcClientConfig
	defaultFabricEP  string

	offlineMode bool
)

func init() {

	viper := config.Viper

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
	grpc["chaincode"] = defaultCCDeployName
	grpc["tlsenabled"] = false
	viper.SetDefault("grpc", grpc)

	// REST Server
	rest := map[string]interface{}{}
	rest["server"] = "example.abchain.org:8080"
	viper.SetDefault("rest", rest)

	// Setting
	setting := map[string]interface{}{}
	setting["offline"] = false
	viper.SetDefault("setting", setting)

}

func StartService() {

	err := config.LoadConfig("conf",
		[]string{"src/hyperledger.abchain.org/cases/ae"},
		nil)

	if err != nil {
		logger.Errorf("Init global config failed: %v", err)
		return
	}

	// Init Fabric
	err = initFabric()
	if err != nil {
		logger.Errorf("Init fabric failed: %v", err)
		return
	}

	viper := config.Viper

	// Init Wallet
	walletFile := filepath.Join(config.FabricPeerFS, viper.GetString("wallet.filename"))
	logger.Debugf("Use wallet file: %s", walletFile)
	defaultWallet = wallet.NewWallet(walletFile)
	err = defaultWallet.Load()
	if err != nil {
		logger.Errorf("Load wallet file failed: %v", err)
		return
	}

	// Init gRPC ClientConfig
	defaultRpcConfig = client.NewRPCConfig(viper.GetString("grpc.chaincode"))
	username := viper.GetString("grpc.username")
	if username != "" {
		defaultRpcConfig.SetUser(username)
		defaultRpcConfig.SetAttrs([]string{
			chaincode.RegionAttr,
			chaincode.PrivilegeAttr,
		}, false)
	}

	// Init REST ClientConfig
	defaultFabricEP = viper.GetString("rest.server")
	logger.Debugf("Use fabric peer REST server: %v", defaultFabricEP)

	offlineMode = viper.GetBool("setting.offline")
	if offlineMode {
		logger.Warning("Running offline mode")
	}

	// Start Local HTTP server
	err = startHttpServer(viper.GetString("service.host"), viper.GetInt("service.port"))
	if err != nil {
		logger.Errorf("StartRPCServer failed: %v", err)
		return
	}
}

func StopService() {

	stopHttpServer()

	defaultWallet.Persist()
	defaultRpcConfig.Quit()
}

func initFabric() error {

	viper := config.Viper

	defaultViperSetting := make(map[string]interface{})

	srvAddr := viper.GetString("grpc.server")
	if srvAddr != "" {
		defaultViperSetting[config.FabricRPC_Addr] = srvAddr
	}

	defaultViperSetting[config.Fabric_DataPath] = viper.GetString("wallet.path")

	if viper.GetBool("grpc.tlsenabled") {
		defaultViperSetting[config.FabricRPC_SSL] = true
		defaultViperSetting[config.FabricRPC_SSLCERT] = viper.GetString("grpc.certfile")
	}

	logLvl := viper.GetString("logging.level")
	if logLvl != "" {
		defaultViperSetting[config.Fabric_LogLevel] = logLvl
	}

	return config.InitFabricPeerEx(defaultViperSetting)
}
