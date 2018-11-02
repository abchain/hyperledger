package abchainTx

import (
	"bytes"
	"encoding/asn1"
	"encoding/base64"
	"errors"
	"math/big"

	abcrypto "hyperledger.abchain.org/core/crypto"
	"hyperledger.abchain.org/core/utils"
)

// AddressInterfaceV1 struct
type AddressInterfaceV1 struct{}

// GetPublicKeyHash for AddressInterfaceV1
func (m AddressInterfaceV1) GetPublicKeyHash(pub *abcrypto.PublicKey) ([]byte, error) {

	if pub == nil || pub.Key == nil || pub.Key.X == nil || pub.Key.Y == nil {
		return nil, errors.New("GetPublicKeyHash: input null pointer")
	}

	type StandardPk struct {
		X *big.Int
		Y *big.Int
	}

	stdpk := StandardPk{X: pub.Key.X, Y: pub.Key.Y}

	rawbytes, err := asn1.Marshal(stdpk)

	if err != nil {
		return nil, err
	}

	return utils.SHA256RIPEMD160(rawbytes)
}

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
