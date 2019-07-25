package addrspace

import (
	"bytes"
	"errors"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
	"strings"
)

type externalAddrVerifier struct {
	AddrSpaceConfig
	rt AddressSpace
}

func Verifier(cfg AddrSpaceConfig) *externalAddrVerifier {
	return &externalAddrVerifier{AddrSpaceConfig: cfg}
}

func (v *externalAddrVerifier) Verify(addr *txutil.Address) error {
	prefix, err := v.rt.QueryPrefix()
	if err != nil {
		return err
	}

	if addrh := addr.Hash; len(addrh) <= txutil.ADDRESS_HASH_LEN {
		return errors.New("Invalid address (too short)")
	} else if bytes.Compare(addrh[:len(prefix)], prefix) != 0 {
		return errors.New("Invalid address (not match)")
	}

	return nil

}

var notInterCCCalling = errors.New("Not inter-chaincode calling")

func (v *externalAddrVerifier) PreHandling(stub shim.ChaincodeStubInterface,
	function string, p txutil.Parser) error {

	if !strings.HasPrefix(function, ".") {
		return notInterCCCalling
	}

	v.rt = v.NewTx(stub, p.GetNonce())
	return nil
}
