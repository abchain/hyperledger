package tx

import (
	"fmt"
	"hyperledger.abchain.org/chaincode/shim"
)

type CollectiveTxs map[string]*ChaincodeTx

func NewCollectiveTxs() CollectiveTxs { return CollectiveTxs(make(map[string]*ChaincodeTx)) }

func (s CollectiveTxs) mergeone(in CollectiveTxs) error {

	for k, v := range in {
		if _, ok := s[k]; ok {
			return fmt.Errorf("Method [%s] is collision", k)
		}

		s[k] = v
	}

	return nil

}

func (s CollectiveTxs) Map() map[string]*ChaincodeTx { return map[string]*ChaincodeTx(s) }

func (s CollectiveTxs) Merge(ins ...CollectiveTxs) (CollectiveTxs, error) {

	for _, in := range ins {
		if err := s.mergeone(in); err != nil {
			return s, err
		}
	}

	return s, nil

}

func (s CollectiveTxs) MustMerge(ins ...CollectiveTxs) CollectiveTxs {

	ret, err := s.Merge(ins...)
	if err != nil {
		panic(fmt.Sprintf("Must Merge fail: %s", err))
	}
	return ret

}

//provide a simple chaincode interface ...
func (s CollectiveTxs) Invoke(stub shim.ChaincodeStubInterface, function string, args [][]byte, _ bool) ([]byte, error) {

	h, ok := s[function]
	if !ok {
		return nil, fmt.Errorf("Method %s is not found", function)
	}

	return h.TxCall(stub, function, args)
}
