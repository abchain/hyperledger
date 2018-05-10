package util

import (
	"github.com/abchain/fabric/core/chaincode/shim"
	"github.com/abchain/fabric/core/chaincode/shim/crypto/attr"
	"github.com/golang/protobuf/ptypes/timestamp"
)

var fabric_support_UnifiedTS = true

//timestamp return by stub used package in their own vendor
func GetTimeStamp(stub shim.ChaincodeStubInterface) (ts *timestamp.Timestamp) {

	if fabric_support_UnifiedTS {
		tss, _ := stub.GetTxTimestamp()
		if tss != nil {
			ts = &timestamp.Timestamp{tss.Seconds, tss.Nanos}
		}
	}
	return
}

func GetAttributes(stub shim.ChaincodeStubInterface,
	attributeName string) string {

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
