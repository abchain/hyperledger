package config

import (
	"errors"
	"github.com/abchain/fabric/peerex"
	vieprWrapper "github.com/abchain/fabric/peerex/viper"
	"os"
	"path/filepath"
)

const (
	FabricRPC_Addr    = "service.cliaddress"
	FabricRPC_SSL     = "peer.tls.serviceenabled"
	FabricRPC_SSLCERT = "peer.tls.rootcert.file"
	Fabric_DataPath   = "peer.fileSystemPath"
	Fabric_LogLevel   = "logging_level"
)

var viper = vieprWrapper.New()
var Viper = viper

//a standard routine to read config from a file and feed then into viper
//any path in goProjPath will be append with GOPATH env, and
//addPath is simply added
//localpath (".") is always used
//environment variables is never read
func LoadConfig(fileName string, goProjPath []string, addPath []string) error {

	viper.SetConfigName(fileName)

	viper.AddConfigPath(".")
	gopath := os.Getenv("GOPATH")
	for _, p := range filepath.SplitList(gopath) {

		for _, c := range goProjPath {
			viper.AddConfigPath(filepath.Join(p, c))
		}
	}

	for _, c := range addPath {
		viper.AddConfigPath(c)
	}

	return viper.ReadInConfig()
}

func handleFSDir(fsDir string) error {

	if fsDir != "" && fsDir != "." {
		err := os.MkdirAll(fsDir, 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

var FabricPeerFS string

//simplely init fabric (std output, not use config ...)
//if we need more specified, use GlobalConfig
func InitFabricPeerEx(settings map[string]interface{}) error {

	if settings[Fabric_DataPath] == nil {
		return errors.New("No data path")
	}

	if settings[FabricRPC_Addr] == nil {
		return errors.New("No RPC Address")
	}

	peerConfig := &peerex.GlobalConfig{
		SkipConfigFile: true,
	}

	peerConfig.InitGlobalWrapper(true, settings)

	FabricPeerFS = peerConfig.GetPeerFS()
	//NOTICE: we provide an empty config for peerex and init
	//is always return error (configfile is not found)
	//so we have to omit it
	//the default settings and environment variables should be still read in
	return handleFSDir(FabricPeerFS)
}

type GlobalConfig struct {
	peerex.GlobalConfig
	Settings map[string]interface{}
	LogFile  bool
}

func (g *GlobalConfig) InitFabricPeerEx() error {
	err := g.GlobalConfig.InitGlobalWrapper(g.LogFile, g.Settings)
	if err != nil {
		return err
	}

	FabricPeerFS = g.GetPeerFS()
	return handleFSDir(FabricPeerFS)
}
