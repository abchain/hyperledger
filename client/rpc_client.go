package client

import (
	"errors"
	"github.com/spf13/viper"
	"hyperledger.abchain.org/chaincode/lib/caller"
	yafabric_cli "hyperledger.abchain.org/client/yafabric"
)

type fabricRPCCfg struct {
	ccName string
	caller func() (rpc.Caller, error)
	yacli  *yafabric_cli.RpcClientConfig
	//TODO: add more implement of clients
}

func NewFabricRPCConfig(ccN string) *fabricRPCCfg {
	return &fabricRPCCfg{ccName: ccN}
}

func (c *fabricRPCCfg) UseYAFabricCli() *yafabric_cli.RpcClientConfig {

	c.yacli = yafabric_cli.NewRPCConfig()
	c.caller = func() (rpc.Caller, error) {
		return c.yacli.GetCaller()
	}

	return c.yacli
}

func (c *fabricRPCCfg) YAFabricCli() *yafabric_cli.RpcClientConfig {
	return c.yacli
}

func (c *fabricRPCCfg) GetCCName() string {
	return c.ccName
}

func (c *fabricRPCCfg) GetCaller() (rpc.Caller, error) {
	if c.caller == nil {
		return nil, errors.New("Not use any client implement")
	}

	return c.caller()
}

func (c *fabricRPCCfg) Quit() {
	c.yacli.Quit()
}
