package msp

import (
	"fmt"
	"os"

	"hyperledger.abchain.org/client/hyfabric/utils"

	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/pkg/errors"
)

//MspEnv 组织信息
type MspEnv struct {
	MspID         string
	MspConfigPath string
	MspType       string //"bccsp", "idemix"  默认bccsp
}

var logger = utils.MustGetLogger("cc_msp")
var localMspType = "bccsp"

func (m *MspEnv) Verify() error {
	if m.MspType == "" {
		m.MspType = localMspType
		logger.Warning("MspType is nll use default type : bccsp")
	}
	m.MspConfigPath = m.MspConfigPath
	if utils.IsNullOrEmpty(m.MspConfigPath) {
		return errors.New("not fund MspConfigPath")
	}
	if utils.IsNullOrEmpty(m.MspType) {
		return errors.New("not fund MspType")
	}

	return nil
}

// InitCrypto initializes crypto for this peer
func (m *MspEnv) InitCrypto() (*factory.FactoryOpts, error) {
	err := m.Verify()
	if err != nil {
		return nil, err
	}
	// Check whether msp folder exists
	fi, err := os.Stat(m.MspConfigPath)
	if os.IsNotExist(err) || !fi.IsDir() {
		// No need to try to load MSP from folder which is not available
		return nil, errors.Errorf("cannot init crypto, missing %s folder", m.MspConfigPath)
	}
	// Check whether localMSPID exists
	if m.MspID == "" {
		return nil, errors.New("the local MSP must have an ID")
	}

	// // Init the BCCSP
	// SetBCCSPKeystorePath()
	// var bccspConfig *factory.FactoryOpts
	//	err = viperutil.EnhancedExactUnmarshalKey("peer.BCCSP", &bccspConfig)

	bccspConfig := NewBccspConf()
	if err != nil {
		return nil, errors.WithMessage(err, "could not parse YAML config")
	}

	err = LoadLocalMspWithType(m.MspConfigPath, bccspConfig, m.MspID, m.MspType)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("error when setting up MSP of type %s from directory %s", m.MspType, m.MspConfigPath))
	}

	return bccspConfig, nil
}

// SetBCCSPKeystorePath sets the file keystore path for the SW BCCSP provider
// to an absolute path relative to the config file
// func SetBCCSPKeystorePath() {
// 	if str, b := utils.GetReplaceKeyAbsPath("peer.BCCSP.SW.FileKeyStore.KeyStore", "peer.mspConfigPath"); b {
// 		viper.Set("peer.BCCSP.SW.FileKeyStore.KeyStore", str)
// 	} else {
// 		viper.Set("peer.BCCSP.SW.FileKeyStore.KeyStore", filepath.Join(str, "KeyStore"))
// 	}
// }

//NewBccspConf 默认，暂时写死，也可以从配置文件中取
func NewBccspConf() *factory.FactoryOpts {
	return &factory.FactoryOpts{
		ProviderName: "SW",
		SwOpts: &factory.SwOpts{
			HashFamily: "SHA2",
			SecLevel:   256,
			Ephemeral:  true,
		},
	}
}
