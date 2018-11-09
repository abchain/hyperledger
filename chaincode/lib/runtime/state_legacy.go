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
	path string
	stub shim.ChaincodeStubInterface
}

type shimStateMapROLegacy struct {
	shimStateMapLegacy
}

func (w *shimStateMapLegacy) SubMap(node string) StateMap_Legacy {

	if node == "" {
		return w
	}

	if node[len(node)-1] != '/' {
		node = node + "/"
	}

	return &shimStateMapLegacy{
		w.path + node,
		w.stub,
	}
}

func (w *shimStateMapLegacy) StoragePath() string {
	return w.path
}

func (w *shimStateMapLegacy) GetRaw(key string) ([]byte, error) {
	return w.stub.GetState(w.path + key)
}

func (w *shimStateMapLegacy) SetRaw(key string, raw []byte) error {
	return w.stub.PutState(w.path+key, raw)
}

func (w *shimStateMapLegacy) Get(key string, m p.Message) error {

	raw, err := w.GetRaw(key)
	if err != nil {
		return err
	}

	if raw == nil {
		return nil
	}

	return p.Unmarshal(raw, m)
}

func (w *shimStateMapLegacy) Set(key string, m p.Message) error {

	raw, err := p.Marshal(m)
	if err != nil {
		return err
	}

	return w.SetRaw(key, raw)
}

func (w *shimStateMapROLegacy) Set(string, p.Message) error {
	return nil
}

//default
func NewShimMapLegacy(root string, stub shim.ChaincodeStubInterface, readOnly bool) StateMap_Legacy {
	if readOnly {
		return &shimStateMapROLegacy{
			shimStateMapLegacy{
				"/" + root + "/",
				stub,
			},
		}
	}

	return &shimStateMapLegacy{
		"/" + root + "/",
		stub,
	}

}
