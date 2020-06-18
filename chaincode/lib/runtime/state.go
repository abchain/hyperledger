package runtime

import (
	"hyperledger.abchain.org/chaincode/shim"
)

type StorageObject interface {
	GetObject() interface{}
	Save() interface{}
	Load(interface{}) error
}

//allowing one or more "extended object" can be persisted
//When load return this error, the rest bytes is unmarshalled
//with the object in Ext field, and Load will be called again
//NOTICE: if there is no any rest bytes, Load is just called
//with the provided Ext field (without any change), and any
//failure in Unmarshalling just caused the whole Get process
//return failure
//NOTICE: Load return more ExtendedObject when there is
//no more bytes left will cause Get returun ExtendedObject as
//an error
type ExtendedObject struct {
	Ext interface{}
}

func (ExtendedObject) Error() string { return "More object" }

type wrappingObject struct {
	obj interface{}
}

func WrapObject(v interface{}) wrappingObject { return wrappingObject{v} }

func (w wrappingObject) GetObject() interface{} { return w.obj }
func (w wrappingObject) Save() interface{}      { return w.obj }
func (w wrappingObject) Load(interface{}) error { return nil }

type StateMap interface {
	SubMap(string) StateMap
	StoragePath() string
	GetRaw(string) ([]byte, error)
	SetRaw(string, []byte) error
	Get(string, StorageObject) error
	Set(string, StorageObject) error
	Delete(string) error
}

//NewShimMap create default shimMap
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

//NewShimMapWithCache create shimMap with a wrapping stub along with write-cache
//which provide a more consistent data-isolation (invoked chaincode read uncommited
//data) on different chain-platform but tamper a little so it is not made an default
//option
func NewShimMapWithCache(root string, stub shim.ChaincodeStubInterface) StateMap {
	return NewShimMap(root, NewCachingStub(stub), false)
}
