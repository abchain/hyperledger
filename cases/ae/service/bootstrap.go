package service

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"hyperledger.abchain.org/applications/asset/wallet"
	simplewallet "hyperledger.abchain.org/applications/asset/wallet/simple"
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

var (
	defaultWallet    wallet.Wallet
	defaultRpcCaller func() (rpc.Caller, error)
	defaultChain     func() (client.ChainInfo, error)

	offlineMode bool
)

const (
	conf_wallet  = "wallet"
	conf_service = "service"
	conf_grpc    = "grpc"
	conf_rest    = "rest"
)

func StartService() {

	if err := config.LoadConfig("conf", []string{"src/hyperledger.abchain.org/cases/ae"}); err != nil {
		logger.Errorf("Init global config failed: %v", err)
		return
	}
	config.LoggingInit(`%{color}%{time:15:04:05.000} %{level:.4s} [%{module:.6s}] %{shortfile} %{shortfunc} ▶ %{message}%{color:reset}`, os.Stdout)

	// Init Wallet
	defaultWallet = simplewallet.LoadWallet(viper.Sub(conf_wallet))
	if err := defaultWallet.Load(); err != nil {
		logger.Errorf("Load wallet file failed: %v", err)
		return
	}
	defer defaultWallet.Persist()

	cfg := client.NewFabricRPCConfig(chaincode.CC_NAME)
	defaultRpcCaller = cfg.GetCaller
	defaultChain = cfg.GetChain
	var err error

	offlineMode = viper.GetBool("offline")
	fmt.Println("get offline:", offlineMode)
	if offlineMode {
		logger.Warning("Running offline mode")
		err = cfg.UseLocalCli()
	} else {
		// Init gRPC ClientConfig
		rpcsetting := viper.Sub(conf_grpc)
		fabricType := rpcsetting.GetString("fabric")
		switch fabricType {
		case "1.x":
			// panic("No implement")
			logger.Debug("use 1.x fabric")
			err = cfg.UseHyFabricCli(rpcsetting)
		case "0.6":
			panic("No implement")
		default:
			cfg.UseYAFabricCli(rpcsetting)
			if viper.IsSet(conf_rest) {
				logger.Infof("Use REST setting [%v] for client", viper.GetStringMap(conf_rest))
				cfg.UseYAFabricREST(viper.Sub(conf_rest))
			}
		}
	}

	if err != nil {
		logger.Errorf("Use client fail: %v", err)
		return
	}

	defer cfg.Quit()

	// start server
	apputil.StartHttpServer(viper.Sub(conf_service), buildRouter())

}

func StopService() {

	apputil.StopHttpServer()
}
