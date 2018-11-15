package client

import (
	"errors"
	"github.com/spf13/viper"
	"hyperledger.abchain.org/chaincode/lib/caller"
)

type RpcSpec struct {
	//notice chaincode name is different to the ccname in txgenerator, the later
	//is used in the hyperledger-project compatible tx
	ChaincodeName string
	Attributes    []string
	Options       *viper.Viper
}

type RpcClient interface {
	Caller(*RpcSpec) (rpc.Caller, error)
	Load(*viper.Viper) error
	Quit()
}

var Client_Impls map[string]func() RpcClient

type fabricRPCCfg struct {
	ccName string
	cli    RpcClient
}

func NewFabricRPCConfig(ccN string) *fabricRPCCfg {
	return &fabricRPCCfg{ccName: ccN}
}

func (c *fabricRPCCfg) UseCli(name string, vp *viper.Viper) error {
	cfg, ok := Client_Impls[name]
	if !ok {
		return errors.New("No implement")
	}

	c.cli = cfg()
	return c.cli.Load(vp)
}

func (c *fabricRPCCfg) UseYAFabricCli(vp *viper.Viper) error {

	return c.UseCli("yafabric", vp)
}

func (c *fabricRPCCfg) GetCCName() string {
	return c.ccName
}

func (c *fabricRPCCfg) GetCaller() (rpc.Caller, error) {
	return c.GetCallerWithSpec(nil)
}

func (c *fabricRPCCfg) GetCallerWithSpec(spec *RpcSpec) (rpc.Caller, error) {
	if c.cli == nil {
		return nil, errors.New("Not use any client implement")
	}

	return c.cli.Caller(spec)
}

func (c *fabricRPCCfg) Quit() {
	if c.cli != nil {
		c.cli.Quit()
	}
	c.cli = nil
}
