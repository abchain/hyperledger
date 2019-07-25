package addrspace

import (
	"encoding/binary"
	"errors"
	"hyperledger.abchain.org/chaincode/impl"
	"hyperledger.abchain.org/chaincode/lib/runtime"
)

type addressSpaceImpl struct {
	*runtime.ChaincodeRuntime
}

func (addrImpl *addressSpaceImpl) getCCName() (string, error) {
	ivf, err := impl.GetInnerInvoke(addrImpl.Stub())
	if err != nil {
		return "", err
	}

	return ivf.GetCallingChaincodeName(), nil
}

func (impl *addressSpaceImpl) RegisterCC() error {

	ccN, err := impl.getCCName()
	if err != nil {
		return err
	}

	var prefix []byte
	err = impl.Storage.Get(ccN, runtime.WrapObject(&prefix))
	if err != nil {
		return err
	} else if prefix != nil {
		return errors.New("Duplicated registering")
	}

	var assignedSeries int
	if err := impl.Storage.Get("@", runtime.WrapObject(&assignedSeries)); err != nil {
		return err
	} else if assignedSeries == 0 {
		//we just inited from a magic code (69)
		assignedSeries = 69
	} else {
		assignedSeries = assignedSeries + 1
	}
	if err := impl.Storage.Set("@", runtime.WrapObject(assignedSeries)); err != nil {
		return err
	}

	buf := make([]byte, binary.MaxVarintLen64)
	buf = buf[:binary.PutUvarint(buf, uint64(assignedSeries))]

	return impl.Storage.Set(ccN, runtime.WrapObject(buf))
}

func (impl *addressSpaceImpl) QueryPrefix() ([]byte, error) {

	ccN, err := impl.getCCName()
	if err != nil {
		return nil, err
	}

	var prefix []byte
	err = impl.Storage.Get(ccN, runtime.WrapObject(&prefix))
	if err != nil {
		return nil, err
	} else if prefix == nil {
		return nil, errors.New("No registry")
	}

	return prefix, nil
}

func (impl *addressSpaceImpl) NormalizeAddress(addr []byte) ([]byte, error) {

	prefix, err := impl.QueryPrefix()
	if err != nil {
		return nil, err
	}

	return append(prefix, addr...), nil

}

type internalAddrSpaceImpl struct{}

func (internalAddrSpaceImpl) RegisterCC() error                            { return nil }
func (internalAddrSpaceImpl) QueryPrefix() ([]byte, error)                 { return []byte{}, nil }
func (internalAddrSpaceImpl) NormalizeAddress(addr []byte) ([]byte, error) { return addr, nil }
