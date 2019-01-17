package abchainTx

import (
	"bytes"
	"encoding/base64"
	"errors"
)

// AddressInterfaceV1 struct
type AddressInterfaceV1 struct{}

// NewAddressFromString for AddressInterfaceV1
func (m AddressInterfaceV1) NewAddressFromString(addressStr string) (*Address, error) {

	data, err := base64.RawURLEncoding.DecodeString(addressStr)

	if err != nil {
		return nil, err
	}

	if len(data) < AddressRequirePartByteSize+AddressVerifyCodeSize {
		return nil, errors.New("Invalid address size")
	}

	if data[0] != networkId {
		return nil, errors.New("Not current network")
	}

	ckborder := len(data) - AddressVerifyCodeSize

	ck, err := GetAddrCheckSum(data[:ckborder])
	if err != nil {
		return nil, err
	}

	if bytes.Compare(ck[:], data[ckborder:]) != 0 {
		return nil, errors.New("checksum error")
	}

	return &Address{
		ADDRESS_VERSION,
		data[0],
		data[1:ckborder],
	}, nil
}

// Serialize for AddressInterfaceV1
func (m AddressInterfaceV1) Serialize(addr *Address) []byte {

	fullbytes := bytes.Join([][]byte{[]byte{addr.NetworkId}, addr.Hash}, nil)

	if len(fullbytes) < AddressRequirePartByteSize {
		return nil
	}

	ck, err := GetAddrCheckSum(fullbytes)
	if err != nil {
		return nil
	}

	return append(fullbytes, ck[:]...)
}
