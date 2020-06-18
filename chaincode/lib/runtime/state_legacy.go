package runtime

import (
	p "github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/shim"
)

type StateMap_Legacy interface {
	SubMap(string) StateMap_Legacy
	StoragePath() string
	GetRaw(string) ([]byte, error)
	SetRaw(string, []byte) error
	Get(string, p.Message) error
	Set(string, p.Message) error
}

type shimStateMapLegacy struct {
	StateMap
}

func (w shimStateMapLegacy) SubMap(node string) StateMap_Legacy {

	return shimStateMapLegacy{w.StateMap.SubMap(node)}
}

func (w shimStateMapLegacy) Get(key string, m p.Message) error {

	raw, err := w.GetRaw(key)
	if err != nil {
		return err
	}

	if raw == nil {
		return nil
	}

	return p.Unmarshal(raw, m)
}

func (w shimStateMapLegacy) Set(key string, m p.Message) error {

	raw, err := p.Marshal(m)
	if err != nil {
		return err
	}

	return w.SetRaw(key, raw)
}

//NewShimMapLegacy create legacy interface
func NewShimMapLegacy(root string, stub shim.ChaincodeStubInterface, readOnly bool) StateMap_Legacy {
	return shimStateMapLegacy{NewShimMap(root, stub, readOnly)}

}
