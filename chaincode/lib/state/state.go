package state

import (
	abchain_shim "github.com/abchain/fabric/core/chaincode/shim"
	p "github.com/golang/protobuf/proto"
)

type StateMap interface {
	SubMap(string) StateMap
	StoragePath() string
	GetRaw(string) ([]byte, error)
	SetRaw(string, []byte) error
	Get(string, p.Message) error
	Set(string, p.Message) error
}

func NewFabric06ShimMap(root string, stub abchain_shim.ChaincodeStubInterface, readOnly bool) StateMap {

	if readOnly {
		return &shimStateMapRO{
			shimStateMap{
				"/" + root + "/",
				stub,
			},
		}
	}

	return &shimStateMap{
		"/" + root + "/",
		stub,
	}

}

//default
func NewShimMap(root string, stub interface{}, readOnly bool) StateMap {
	switch s := stub.(type) {
	case abchain_shim.ChaincodeStubInterface:
		return NewFabric06ShimMap(root, s, readOnly)
	default:
		panic("No corresponding stub")
	}

}
