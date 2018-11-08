package runtime

import (
	"hyperledger.abchain.org/chaincode/shim"
)

type StorageObject interface {
	GetObject() interface{}
	Save() interface{}
	Load(interface{}) error
}

type StateMap interface {
	SubMap(string) StateMap
	StoragePath() string
	GetRaw(string) ([]byte, error)
	SetRaw(string, []byte) error
	Get(string, StorageObject) error
	Set(string, StorageObject) error
}

//default
func NewShimMap(root string, stub shim.ChaincodeStubInterface, readOnly bool) StateMap {
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
