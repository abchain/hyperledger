package addrspace

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"sync"
)

type GeneralCall struct {
	txgen.TxCaller
}

const (
	Method_Reg   = "ADDRSPACE.REG"
	Method_Query = "ADDRSPACE.QUERY"
)

func (i GeneralCall) RegisterCC() error {

	return i.Invoke(Method_Reg, new(empty.Empty))
}

func (i GeneralCall) QueryPrefix() ([]byte, error) {

	ret, err := i.Query(Method_Query, new(empty.Empty))

	if err != nil {
		return nil, err
	}

	p := new(wrappers.BytesValue)

	err = txgen.SyncQueryResult(p, ret)
	if err != nil {
		return nil, err
	}

	return p.GetValue(), nil
}

//notice: there is not really a NormalizeAddress call
//(we just do a query and use the returning result)
func (i GeneralCall) NormalizeAddress(addr []byte) ([]byte, error) {

	prefix, err := i.QueryPrefix()
	if err != nil {
		return nil, err
	}

	return append(prefix, addr...), nil
}

func (i GeneralCall) CacheImpl(c *generalCallCache) cachedGeneralCall {
	return cachedGeneralCall{c, i}
}

type generalCallCache struct {
	sync.RWMutex
	prefixCache []byte
}

func CallCache() *generalCallCache { return new(generalCallCache) }

type cachedGeneralCall struct {
	*generalCallCache
	GeneralCall
}

func (i cachedGeneralCall) QueryPrefix() ([]byte, error) {

	i.RLock()
	if i.prefixCache != nil {
		i.RUnlock()
		return i.prefixCache, nil
	}
	i.RUnlock()

	ret, err := i.GeneralCall.QueryPrefix()

	if err != nil {
		return nil, err
	}

	i.Lock()
	i.prefixCache = ret
	i.Unlock()

	return ret, nil
}

func (i cachedGeneralCall) NormalizeAddress(addr []byte) ([]byte, error) {

	prefix, err := i.QueryPrefix()
	if err != nil {
		return nil, err
	}

	return append(prefix, addr...), nil
}
