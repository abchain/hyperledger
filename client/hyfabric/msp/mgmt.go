/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msp

import (
	"sync"

	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/msp/cache"
	"github.com/pkg/errors"
)

// FIXME: AS SOON AS THE CHAIN MANAGEMENT CODE IS COMPLETE,
// THESE MAPS AND HELPSER FUNCTIONS SHOULD DISAPPEAR BECAUSE
// OWNERSHIP OF PER-CHAIN MSP MANAGERS WILL BE HANDLED BY IT;
// HOWEVER IN THE INTERIM, THESE HELPER FUNCTIONS ARE REQUIRED

var (
	m        sync.Mutex
	localMsp msp.MSP

	localMsptype string
	// var mspMap map[string]msp.MSPManager = make(map[string]msp.MSPManager)
)

// LoadLocalMspWithType loads the local MSP with the specified type from the specified directory
func LoadLocalMspWithType(dir string, bccspConfig *factory.FactoryOpts, mspID, mspType string) error {
	if mspID == "" {
		return errors.New("the local MSP must have an ID")
	}
	//github.com\hyperledger\fabric\msp\configbuilder.go 将msp路径下的各个文件读取 之后proto序列化
	conf, err := msp.GetLocalMspConfigWithType(dir, bccspConfig, mspID, mspType)
	if err != nil {
		return err
	}
	localMsptype = mspType
	//Setup  实现：github.com\hyperledger\fabric\msp\mspimplsetup.go
	//将conf用proto反序列化之后为github.com\hyperledger\fabric\protos\msp\msp_config.pb.go中的FabricMSPConfig的每一项
	//之后通过setup构造成\github.com\hyperledger\fabric\msp\mspimpl.go中的bccspmsp
	return GetLocalMSP().Setup(conf)
}

// GetLocalMSP returns the local msp (and creates it if it doesn't exist)
func GetLocalMSP() msp.MSP {
	if localMsp != nil {
		return localMsp
	}
	m.Lock()
	defer m.Unlock()

	if localMsp != nil {
		return localMsp
	}

	localMsp = loadLocaMSP()

	return localMsp
}

func loadLocaMSP() msp.MSP {
	// determine the type of MSP (by default, we'll use bccspMSP)
	// mspType = viper.GetString("peer.localMspType")
	mspType := localMsptype
	if mspType == "" {
		mspType = msp.ProviderTypeToString(msp.FABRIC)
	}

	var mspOpts = map[string]msp.NewOpts{
		msp.ProviderTypeToString(msp.FABRIC): &msp.BCCSPNewOpts{NewBaseOpts: msp.NewBaseOpts{Version: msp.MSPv1_0}},
		msp.ProviderTypeToString(msp.IDEMIX): &msp.IdemixNewOpts{msp.NewBaseOpts{Version: msp.MSPv1_1}},
	}
	newOpts, found := mspOpts[mspType]
	if !found {
		logger.Panicf("msp type " + mspType + " unknown")
	}

	mspInst, err := msp.New(newOpts)
	if err != nil {
		logger.Fatalf("Failed to initialize local MSP, received err %+v", err)
	}
	switch mspType {
	case msp.ProviderTypeToString(msp.FABRIC):
		mspInst, err = cache.New(mspInst)
		if err != nil {
			logger.Fatalf("Failed to initialize local MSP, received err %+v", err)
		}
	case msp.ProviderTypeToString(msp.IDEMIX):
		// Do nothing
	default:
		panic("msp type " + mspType + " unknown")
	}

	logger.Debug("Created new local MSP")

	return mspInst
}

//GetSigningIdentity 获取bccsp类型的签名身份
func GetSigningIdentity() (msp.SigningIdentity, error) {
	return GetLocalMSP().GetDefaultSigningIdentity()
}

// //GetDefaultSigner 获取默认签名
// func GetDefaultSigner() (msp.SigningIdentity, error) {
// 	signer, err := mspmgmt.GetSigningIdentity()
// 	if err != nil {
// 		return nil, errors.WithMessage(err, "error obtaining the default signing identity")
// 	}

// 	return signer, err
// }
