package service

import (
	"os"

	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"hyperledger.abchain.org/applications/asset/wallet"
	apputil "hyperledger.abchain.org/applications/util"
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
	"hyperledger.abchain.org/chaincode/lib/caller"
	"hyperledger.abchain.org/client"
	"hyperledger.abchain.org/core/config"
)

const (
	defaultCCDeployName = "aecc"
)

var logger = logging.MustGetLogger("server")

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
	defaultRpcConfig apputil.FabricRPCCfg
	defaultFabricEP  string

	offlineMode bool
)

const (
	conf_wallet  = "wallet"
	conf_service = "service"
	conf_grpc    = "grpc"
)

func StartService() {

	config.LoggingInit(`%{color}%{time:15:04:05.000} %{level:.4s} [%{module:.6s}] %{shortfile} %{shortfunc} ▶ %{message}%{color:reset}`, os.Stdout)

	if err := config.LoadConfig("conf", []string{"src/hyperledger.abchain.org/cases/ae"}); err != nil {
		logger.Errorf("Init global config failed: %v", err)
		return
	}

	// Init Wallet
	defaultWallet = wallet.LoadWallet(viper.Sub(conf_wallet))
	if err := defaultWallet.Load(); err != nil {
		logger.Errorf("Load wallet file failed: %v", err)
		return
	}

	offlineMode = viper.GetBool("offline")

	if offlineMode {
		logger.Warning("Running offline mode")
		defaultRpcConfig = rpcCfg(chaincode.CC_NAME)
	} else {
		// Init gRPC ClientConfig
		cfg := client.NewFabricRPCConfig(chaincode.CC_NAME)
		defaultRpcConfig = cfg

		rpcsetting := viper.Sub(conf_grpc)
		fabricType := rpcsetting.GetString("fabric")

		switch fabricType {
		case "1.x":
			panic("No implement")
		case "0.6":
			panic("No implement")
		default:
			cfg.UseYAFabricCli(rpcsetting)
		}
	}

	// Init REST ClientConfig
	defaultFabricEP = viper.GetString("rest.server")
	logger.Debugf("Use fabric peer REST server: %v", defaultFabricEP)

	// start server
	apputil.StartHttpServer(viper.Sub(conf_service), buildRouter())

}

func StopService() {

	apputil.StopHttpServer()

	if defaultWallet != nil {
		defaultWallet.Persist()
	}

	defaultRpcConfig.Quit()
}
