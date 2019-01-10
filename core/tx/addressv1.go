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

	if len(data) != AddressFullByteSize {
		return nil, errors.New("Invalid address size")
	}

	if data[0] != networkId {
		return nil, errors.New("Not current network")
	}

	ck, err := GetAddrCheckSum(data[:AddressPartByteSize])
	if err != nil {
		return nil, err
	}

	if bytes.Compare(ck[:], data[AddressPartByteSize:]) != 0 {
		return nil, errors.New("checksum error")
	}

	return &Address{
		ADDRESS_VERSION,
		data[0],
		data[1:AddressPartByteSize],
	}, nil
}

// Serialize for AddressInterfaceV1
func (m AddressInterfaceV1) Serialize(addr *Address) []byte {

	fullbytes := bytes.Join([][]byte{[]byte{addr.NetworkId}, addr.Hash}, nil)

	if len(fullbytes) != AddressPartByteSize {
		return nil
	}

	ck, err := GetAddrCheckSum(fullbytes)
	if err != nil {
		return nil
	}

	return append(fullbytes, ck[:]...)
}
