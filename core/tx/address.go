package abchainTx

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	abcrypto "hyperledger.abchain.org/core/crypto"
	"hyperledger.abchain.org/core/utils"
	pb "hyperledger.abchain.org/protos"
)

const (
	ADDRESS_VERSION  = 1
	ADDRESS_CRC_LEN  = 4
	ADDRESS_HASH_LEN = 20
)

type Address struct {
	Version   int32
	NetworkId uint8
	Hash      []byte
}

// AddressInterface interface
type AddressInterface interface {
	NewAddressFromString(addressStr string) (*Address, error)
	Serialize(addr *Address) []byte
}

var addrimplv1 AddressInterface = AddressInterfaceV1{}
var addrimpl = addrimplv1

type AddressHelper map[uint8]string

var AddrHelper AddressHelper = map[uint8]string{
	1: "ABCHAIN",
	4: "VCTMAINNET",
}

var networkId uint8 = 1

func init() {
	pb.TxAddrMarshaller = txAddrMarshaller
}

func getNetwork(nId uint8) string {
	n, ok := AddrHelper[nId]
	if !ok {
		return "UNN"
	}

	return n
}

func UseNetwork(i uint8) {

	_, existed := AddrHelper[i]
	if !existed {
		panic(fmt.Sprintf("Network id %d is not register yet", i))
	}

	networkId = i
}

func DefaultNetworkName() string { return getNetwork(networkId) }

// SetAddressInterfaceImpl set impl nil to restore v1
func SetAddressInterfaceImpl(impl AddressInterface) {
	if impl == nil {
		addrimpl = addrimplv1
	} else {
		addrimpl = impl
	}
}

func NormalizeHash(h []byte) []byte {
	if len(h) < ADDRESS_HASH_LEN {
		h = bytes.Join([][]byte{h, make([]byte, ADDRESS_HASH_LEN)}, nil)
	}

	return h[:ADDRESS_HASH_LEN]
}

func NewAddressFromHash(h []byte) *Address {

	if len(h) < ADDRESS_HASH_LEN {
		h = bytes.Join([][]byte{h, make([]byte, ADDRESS_HASH_LEN)}, nil)
	}

	return &Address{
		ADDRESS_VERSION,
		networkId,
		h,
	}
}

func NewAddressFromPrivateKey(priv abcrypto.Signer) (*Address, error) {

	if priv == nil {
		return nil, errors.New("AddressFromPrivateKey: input null pointer")
	}

	return NewAddress(priv.Public())
}

func NewAddress(pub abcrypto.Verifier) (*Address, error) {

	if pub == nil {
		return nil, errors.New("AddressFromPublickKey: input null pointer")
	}

	return NewAddressFromHash(NormalizeHash(pub.Digest())), nil
}

func NewAddressFromPBMessage(addrProto *pb.TxAddr) (*Address, error) {

	if addrProto == nil {
		return nil, errors.New("AddressFromPBMessage: input null pointer")
	}

	return &Address{
		ADDRESS_VERSION,
		networkId,
		addrProto.Hash,
	}, nil
}

const (
	AddressRequirePartByteSize = ADDRESS_HASH_LEN + ADDRESS_VERSION
	AddressVerifyCodeSize      = ADDRESS_CRC_LEN
	// AddressFullByteSize   = AddressVerifyCodeSize + AddressPartByteSize
)

func GetAddrCheckSum(rb []byte) ([AddressVerifyCodeSize]byte, error) {

	hash, err := utils.DoubleSHA256(rb)
	if err != nil {
		return [AddressVerifyCodeSize]byte{}, err
	}

	return [AddressVerifyCodeSize]byte{hash[0], hash[1], hash[2], hash[3]}, nil
}

func NewAddressFromString(addressStr string) (*Address, error) {

	return addrimpl.NewAddressFromString(addressStr)
}

func (addr *Address) Internal() []byte {

	if addr == nil {
		return nil
	}

	return addr.Hash
}

func (addr *Address) PBMessage() *pb.TxAddr {

	if addr == nil {
		return new(pb.TxAddr)
	}

	addrProto := &pb.TxAddr{
		Hash: addr.Hash,
	}

	return addrProto
}

func (addr *Address) Serialize() []byte {

	return addrimpl.Serialize(addr)
}

func txAddrMarshaller(m *pb.TxAddr) ([]byte, error) {

	if m.GetHash() == nil {
		return json.Marshal("NO ADDRESS")
	}

	caddr := NewAddressFromHash(m.GetHash())
	return json.Marshal(caddr.ToString())
}

func (addr *Address) ToString() string {

	return base64.RawURLEncoding.EncodeToString(addr.Serialize())
}

func (addr *Address) Dump() string {

	return fmt.Sprintf("&{Version: %d, Network: %v, Hash: %v}",
		addr.Version, addr.NetworkId, addr.Hash)
}

func (addr *Address) IsExternal() bool {

	return len(addr.Hash) > ADDRESS_HASH_LEN

}

func (addr *Address) IsEqual(otherAddr *Address) bool {

	if otherAddr == nil {
		return false
	}

	if !bytes.Equal(addr.Hash, otherAddr.Hash) {
		return false
	}

	return true
}
