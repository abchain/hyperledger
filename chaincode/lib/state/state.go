package state

import (
	p "github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/shim"
)

type StateMap interface {
	SubMap(string) StateMap
	StoragePath() string
	GetRaw(string) ([]byte, error)
	SetRaw(string, []byte) error
	Get(string, p.Message) error
	Set(string, p.Message) error
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
