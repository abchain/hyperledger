package util

import (
	"github.com/abchain/fabric/core/chaincode/shim"
	"github.com/abchain/fabric/core/chaincode/shim/crypto/attr"
)

func getF06Attributes(stub shim.ChaincodeStubInterface, attributeName string) string {

	attrHandler, err := attr.NewAttributesHandlerImpl(stub)
	if err != nil {
		return ""
	}

	var attrStr string
	attr, err := attrHandler.GetValue(attributeName)
	if err != nil {
		attrStr = ""
	} else {
		attrStr = string(attr)
	}

	return attrStr
}
