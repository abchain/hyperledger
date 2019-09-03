package runtime

import (
	"encoding/asn1"
	"errors"
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
	if key == "" {
		return errors.New("Empty key is not allowed")
	}
	return w.stub.PutState(w.path+key, raw)
}

func (w *shimStateMap) Delete(key string) error {
	return w.stub.DelState(w.path + key)
}

func (w *shimStateMap) Get(key string, m StorageObject) error {

	raw, err := w.GetRaw(key)
	if err != nil {
		return err
	}

	if raw == nil {
		return nil
	}
	obj := m.GetObject()
	rest, err := asn1.Unmarshal(raw, obj)
	if err != nil {
		return err
	}
	err = m.Load(obj)
	for extObj, ok := err.(ExtendedObject); ok; extObj, ok = err.(ExtendedObject) {
		if len(rest) == 0 {
			return m.Load(extObj.Ext)
		}
		obj = extObj.Ext
		rest, err = asn1.Unmarshal(rest, obj)
		if err != nil {
			return err
		}
		err = m.Load(obj)
	}

	return err
}

func SeralizeObject(m StorageObject) ([]byte, error) {
	return asn1.Marshal(m.Save())
}

func (w *shimStateMap) Set(key string, m StorageObject) error {

	if key == "" {
		return errors.New("Empty key is not allowed")
	}

	raw, err := SeralizeObject(m)
	if err != nil {
		return err
	}

	return w.SetRaw(key, raw)
}

func (w *shimStateMapRO) Set(key string, _ StorageObject) error {
	if key == "" {
		return errors.New("Empty key is not allowed")
	}
	return nil
}

func (w *shimStateMapRO) SubMap(node string) StateMap {

	s := w.shimStateMap.SubMap(node).(*shimStateMap)
	return &shimStateMapRO{*s}

}

// func (w *shimStateMap) Get(key string, m p.Message) error {

// 	raw, err := w.GetRaw(key)
// 	if err != nil {
// 		return err
// 	}

// 	if raw == nil {
// 		return nil
// 	}

// 	return p.Unmarshal(raw, m)
// }

// func (w *shimStateMap) Set(key string, m p.Message) error {

// 	raw, err := p.Marshal(m)
// 	if err != nil {
// 		return err
// 	}

// 	return w.SetRaw(key, raw)
// }

// func (w *shimStateMapRO) Set(string, p.Message) error {
// 	return nil
// }

func (w *shimStateMapRO) SetRaw(string, []byte) error {
	return nil
}
