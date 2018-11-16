package service

import (
	"errors"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"hyperledger.abchain.org/applications/asset/wallet"
	"hyperledger.abchain.org/applications/blockchain"
	apputil "hyperledger.abchain.org/applications/util"
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
	"hyperledger.abchain.org/chaincode/lib/caller"
	"hyperledger.abchain.org/client"
	"hyperledger.abchain.org/core/config"
	"os"
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

func (rpcCfg) GetChain() (client.ChainInfo, error) {
	return nil, errors.New("Not support")
}

var (
	defaultWallet      wallet.Wallet
	defaultRpcConfig   apputil.FabricRPCCfg
	defaultChainConfig blockchain.FabricChainCfg

	offlineMode bool
)

const (
	conf_wallet  = "wallet"
	conf_service = "service"
	conf_grpc    = "grpc"
	conf_rest    = "rest"
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
		cfg := rpcCfg(chaincode.CC_NAME)
		defaultRpcConfig = cfg
		defaultChainConfig = cfg
	} else {
		// Init gRPC ClientConfig
		cfg := client.NewFabricRPCConfig(chaincode.CC_NAME)
		defaultRpcConfig = cfg
		defaultChainConfig = cfg

		rpcsetting := viper.Sub(conf_grpc)
		fabricType := rpcsetting.GetString("fabric")

		switch fabricType {
		case "1.x":
			panic("No implement")
		case "0.6":
			panic("No implement")
		default:
			cfg.UseYAFabricCli(rpcsetting)
			cfg.UseYAFabricREST(viper.Sub(conf_rest))
		}
	}

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
