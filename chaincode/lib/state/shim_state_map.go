package state

import (
	p "github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/shim"
)

type shimStateMap struct {
	path string
	stub shim.ChaincodeStubInterface
}

type shimStateMapRO struct {
	shimStateMap
}

func (w *shimStateMap) SubMap(node string) StateMap {

	if node == "" {
		return w
	}

	if node[len(node)-1] != '/' {
		node = node + "/"
	}

	return &shimStateMap{
		w.path + node,
		w.stub,
	}
}

func (w *shimStateMap) StoragePath() string {
	return w.path
}

func (w *shimStateMap) GetRaw(key string) ([]byte, error) {
	return w.stub.GetState(w.path + key)
}

func (w *shimStateMap) SetRaw(key string, raw []byte) error {
	return w.stub.PutState(w.path+key, raw)
}

func (w *shimStateMap) Get(key string, m p.Message) error {

	raw, err := w.GetRaw(key)
	if err != nil {
		return err
	}

	if raw == nil {
		return nil
	}

	return p.Unmarshal(raw, m)
}

func (w *shimStateMap) Set(key string, m p.Message) error {

	raw, err := p.Marshal(m)
	if err != nil {
		return err
	}

	return w.SetRaw(key, raw)
}

func (w *shimStateMapRO) Set(string, p.Message) error {
	return nil
}

func (w *shimStateMapRO) SetRaw(string, []byte) error {
	return nil
}
