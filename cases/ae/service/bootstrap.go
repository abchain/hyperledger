package service

import (
	"hyperledger.abchain.org/asset/wallet"
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
	"hyperledger.abchain.org/chaincode/lib/caller"
	"hyperledger.abchain.org/client"
	"hyperledger.abchain.org/config"
	"path/filepath"
)

const (
	defaultCCDeployName = "aecc"
)

type rpcCfg string

func (s rpcCfg) GetCCName() string {
	return string(s)
}

func (rpcCfg) Quit() {

}

func (rpcCfg) GetCaller() (rpc.Caller, error) {
	return ccCaller, nil
}

var (
	defaultWallet    wallet.Wallet
	defaultRpcConfig client.FabricRPCCfg
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

	// Init REST ClientConfig
	defaultFabricEP = viper.GetString("rest.server")
	logger.Debugf("Use fabric peer REST server: %v", defaultFabricEP)

	offlineMode = viper.GetBool("setting.offline")
	if offlineMode {
		logger.Warning("Running offline mode")
		defaultRpcConfig = rpcCfg(chaincode.CC_NAME)
	} else {
		// Init gRPC ClientConfig
		cfg := client.NewFabricRPCConfig(chaincode.CC_NAME)

		// TODO: allow different implenents of client
		yacfg := cfg.UseYAFabricCli(viper.GetString("grpc.chaincode"))
		username := viper.GetString("grpc.username")
		if username != "" {
			yacfg.SetUser(username)
			yacfg.SetAttrs([]string{
				chaincode.RegionAttr,
				chaincode.PrivilegeAttr,
			}, false)
		}

		defaultRpcConfig = cfg
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
	defaultViperSetting[config.FabricRPC_SSL] = viper.GetBool("grpc.tlsenabled")
	defaultViperSetting[config.FabricRPC_SSLCERT] = viper.GetString("grpc.certfile")

	logLvl := viper.GetString("logging.level")
	if logLvl != "" {
		defaultViperSetting[config.Fabric_LogLevel] = logLvl
	}

	return config.InitFabricPeerEx(defaultViperSetting)
}
