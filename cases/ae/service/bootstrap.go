package service

import (
	"github.com/abchain/fabric/peerex"
	"hyperledger.abchain.org/asset/wallet"
	"hyperledger.abchain.org/client"
	"hyperledger.abchain.org/config"
	"os"
	"path/filepath"
)

var (
	DefaultWallet     wallet.Wallet
	DefaultRPCClient  *client.RPCClient
	DefaultRESTClient *client.RESTClient

	offlineMode bool
)

func StartService() {
	var err error

	// Read config
	globalConfig := &config.GlobalConfig{ConfigFileName: "conf"}

	globalConfig.ConfigPath = []string{"."}

	gopath := os.Getenv("GOPATH")
	for _, p := range filepath.SplitList(gopath) {
		confpath := filepath.Join(p, "src/hyperledger.abchain.org/cases/ae")
		globalConfig.ConfigPath = append(globalConfig.ConfigPath, confpath)
	}

	err = globalConfig.InitGlobal()
	if err != nil {
		logger.Errorf("Init global failed: %v", err)
		return
	}

	// Init Fabric
	if err = initFabric(globalConfig); err != nil {
		logger.Errorf("Init fabric failed: %v", err)
		return
	}

	viper := config.Viper

	// Init Wallet
	path := viper.GetString("wallet.path")
	filename := viper.GetString("wallet.filename")
	walletFile := filepath.Join(path, filename)
	err = os.MkdirAll(path, 0777)
	if err != nil {
		logger.Errorf("mkdir %v failed: %v", path, err)
		return
	}
	logger.Debugf("Use wallet file: %s", walletFile)
	DefaultWallet = wallet.NewWallet(walletFile)
	err = DefaultWallet.Load()
	if err != nil {
		logger.Errorf("Load wallet file failed: %v", err)
		return
	}

	// Init gRPC Client
	DefaultRPCClient, err = client.NewRPCClient()
	if err != nil {
		logger.Errorf("Create RPC Client failed: %v", err)
		return
	}

	// Init REST Client
	RESTServer := viper.GetString("rest.server")
	logger.Debugf("Connect to REST server: %v", RESTServer)
	DefaultRESTClient, err = client.NewRESTClient(RESTServer)
	if err != nil {
		logger.Errorf("Create REST Client failed: %v", err)
		return
	}

	// Connect to gRPC Server
	offlineMode = viper.GetBool("setting.offline")
	if !offlineMode {
		// Connect gRPC Server
		gRPCServer := viper.GetString("grpc.server")
		err = DefaultRPCClient.Connect(gRPCServer)
		if err != nil {
			logger.Errorf("Connect gRPC server(%v) failed: %v", gRPCServer, err)
			return
		}

		err = DefaultRPCClient.SetSecurityPolicy(viper.GetString("grpc.username")) // Set UserName
		if err != nil {
			logger.Errorf("SetSecurityPolicy failed: %v", err)
			return
		}

		err = DefaultRPCClient.SetChaincodeName(viper.GetString("grpc.chaincode")) // Set Chaincode Name
		if err != nil {
			logger.Errorf("SetChaincodeName failed: %v", err)
			return
		}
	} else {
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

	DefaultWallet.Persist()
	DefaultRPCClient.Close()
}

func initFabric(cfg *config.GlobalConfig) error {
	peerConfig := &peerex.GlobalConfig{}
	peerConfig.ConfigFileName = "conf"
	peerConfig.ConfigPath = cfg.ConfigPath

	viper := config.Viper

	defaultViperSetting := make(map[string]interface{})
	defaultViperSetting["peer.fileSystemPath"] = viper.GetString("grpc.path")
	defaultViperSetting["peer.tls.rootcert.file"] = viper.GetString("grpc.certfile")
	defaultViperSetting["peer.tls.serviceenabled"] = viper.GetBool("grpc.tlsenabled")

	err := peerConfig.InitGlobalWrapper(true, defaultViperSetting)
	if err != nil {
		return err
	}

	var fsDir string = defaultViperSetting["peer.fileSystemPath"].(string)
	if fsDir != "" && fsDir != "." {
		err = os.MkdirAll(viper.GetString("grpc.path"), 0777)
		if err != nil {
			return err
		}
	}

	return nil
}
