package util

import (
	abchain_shim "github.com/abchain/fabric/core/chaincode/shim"
	"github.com/golang/protobuf/ptypes/timestamp"
)

//timestamp return by stub used package in their own vendor
func GetTimeStamp(stub interface{}) (ts *timestamp.Timestamp) {

	switch s := stub.(type) {
	case abchain_shim.ChaincodeStubInterface:
		tss, _ := s.GetTxTimestamp()
		if tss != nil {
			ts = &timestamp.Timestamp{tss.Seconds, tss.Nanos}
		}
		return

	default:
		panic("No corresponding stub")
	}
}

func GetTxID(stub interface{}) string {
	switch s := stub.(type) {
	case abchain_shim.ChaincodeStubInterface:
		return s.GetTxID()

	default:
		panic("No corresponding stub")
	}
}

func GetAttributes(stub interface{}, attributeName string) string {

	switch s := stub.(type) {
	case abchain_shim.ChaincodeStubInterface:
		return getF06Attributes(s, attributeName)

	default:
		panic("No corresponding stub")
	}
}
