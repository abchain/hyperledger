package state

import (
	"github.com/abchain/fabric/core/chaincode/shim"
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

func NewFabric06ShimMap(root string, stub shim.ChaincodeStubInterface, readOnly bool) StateMap {

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
func NewShimMap(root string, stub shim.ChaincodeStubInterface, readOnly bool) StateMap {
	return NewFabric06ShimMap(root, stub, readOnly)
}
