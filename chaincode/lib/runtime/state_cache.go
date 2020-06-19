package runtime

import (
	"hyperledger.abchain.org/chaincode/shim"
)

//CachingSupport enable an implement indicate it has support caching,
//i.e. GetState can read uncommited value
type CachingSupport interface {
	CanReadUnCommit()
}

type cachingCcStub struct {
	shim.ChaincodeStubInterface
	stateCache map[string][]byte
}

//NewCachingStub create an stub with local caching
func NewCachingStub(stub shim.ChaincodeStubInterface) shim.ChaincodeStubInterface {
	return &cachingCcStub{stub, make(map[string][]byte)}
}

func (*cachingCcStub) CanReadUnCommit() {}

func (wc *cachingCcStub) GetState(key string) ([]byte, error) {

	if v, existed := wc.stateCache[key]; existed {
		return v, nil
	}

	return wc.ChaincodeStubInterface.GetState(key)

}

func (wc *cachingCcStub) PutState(key string, value []byte) error {

	if err := wc.ChaincodeStubInterface.PutState(key, value); err != nil {
		return err
	}

	//notice we must deep copy the put value
	deepclone := make([]byte, len(value))
	copy(deepclone, value)

	wc.stateCache[key] = deepclone
	return nil
}
