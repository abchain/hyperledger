package runtime

import (
	"hyperledger.abchain.org/chaincode/shim"
)

type cachingCcStub struct {
	shim.ChaincodeStubInterface
	stateCache map[string][]byte
}

//NewCachingStub create an stub with local caching
func NewCachingStub(stub shim.ChaincodeStubInterface) shim.ChaincodeStubInterface {
	return &cachingCcStub{stub, make(map[string][]byte)}
}

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
