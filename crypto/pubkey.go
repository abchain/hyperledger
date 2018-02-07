package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/protos"
	"hyperledger.abchain.org/utils"
	"math/big"
)

type PublicKey struct {
	Version int32

	// Curve type
	CurveType int32

	// underlying public key
	Key *ecdsa.PublicKey

	// underlying root public key fingerprint
	RootFingerPrint []byte

	// index of child public key
	Index *big.Int

	// chaincode
	Chaincode []byte
}

const (
	SECP256K1 int32 = iota + 1
	ECP256_FIPS186
)

const (
	PUBLICKEY_VERSION = 1

	PUBLICKEY_FINGERPRINT_LEN = 8

	PublicKeyCompressedLength = 33
)

var (
	DefaultCurveType = ECP256_FIPS186

	DefaultPublicKeyVersion int32 = PUBLICKEY_VERSION
)

func GetEC(curveType int32) (elliptic.Curve, error) {
	switch curveType {
	case ECP256_FIPS186:
		return elliptic.P256(), nil
	case SECP256K1:
		return Secp256k1(), nil
	default:
		return nil, fmt.Errorf("%d is not a valid curve defination", curveType)
	}
}

func PublicKeyFromBytes(raw []byte) (*PublicKey, error) {

	pubProto := &protos.PublicKey{}
	err := proto.Unmarshal(raw, pubProto)
	if err != nil {
		return nil, err
	}

	return PublicKeyFromPBMessage(pubProto)
}

func PublicKeyFromPBMessage(pubProto *protos.PublicKey) (*PublicKey, error) {

	if pubProto == nil {
		return nil, errors.New("PublicKeyFromPBMessage: input null pointer")
	}

	if pubProto.Version != PUBLICKEY_VERSION {
		return nil, errors.New("Unknown public key version")
	}

	curve, err := GetEC(pubProto.Curvetype)
	if err != nil {
		return nil, err
	}

	X := big.NewInt(0)
	Y := big.NewInt(0)

	// Compress algorithm only works for SECP256K1
	if pubProto.Curvetype == SECP256K1 {
		X, Y = expandPublicKey(pubProto.GetP().X, curve.Params())
	} else {
		X.SetBytes(pubProto.GetP().X)
		Y.SetBytes(pubProto.GetP().Y)
	}

	ecdsaKey := &ecdsa.PublicKey{
		curve,
		X,
		Y,
	}

	pubKey := &PublicKey{
		DefaultPublicKeyVersion,
		pubProto.Curvetype,
		ecdsaKey,
		pubProto.RootFingerprint,
		big.NewInt(0).SetBytes(pubProto.Index),
		pubProto.Chaincode,
	}

	return pubKey, nil
}

func (pub *PublicKey) PBMessage() *protos.PublicKey {

	var X []byte
	var Y []byte

	// Compress algorithm only works for SECP256K1
	if pub.CurveType == SECP256K1 {
		X = compressPublicKey(pub.Key.X, pub.Key.Y)
		Y = nil
	} else {
		X = pub.Key.X.Bytes()
		Y = pub.Key.Y.Bytes()
	}

	pubProto := &protos.PublicKey{
		pub.Version,
		pub.CurveType,
		pub.RootFingerPrint,
		pub.Index.Bytes(),
		pub.Chaincode,
		&protos.ECPoint{
			X,
			Y,
		},
	}

	return pubProto
}

func (pub *PublicKey) Serialize() []byte {

	pubProto := pub.PBMessage()

	raw, err := proto.Marshal(pubProto)
	if err != nil {
		return nil
	}

	return raw
}

func (pub *PublicKey) Str() string {
	return base64.URLEncoding.EncodeToString(pub.Serialize())
}

func (pub *PublicKey) RootFingerPrintStr() string {

	return base64.URLEncoding.EncodeToString(pub.RootFingerPrint)
}

func (pub *PublicKey) String() string {

	return fmt.Sprintf("&{Version: %v, CurveType: %d, Key.X: %v, Key.Y: %v, "+
		"RootFingerPrint: %v, Index: %v, Chaincode: %v}",
		pub.Version, pub.CurveType, pub.Key.X.Bytes(), pub.Key.Y.Bytes(),
		pub.RootFingerPrint, pub.Index, pub.Chaincode)
}

func (pub *PublicKey) ChildKey(index *big.Int) (*PublicKey, error) {

	if index.Int64() == 0 {
		return pub, nil
	}

	return getChildPublicKey(pub, index)
}

func (pub *PublicKey) Verify(hash []byte, sig *Signature) bool {

	return ecdsa.Verify(pub.Key, hash, sig.R, sig.S)
}

func (pub *PublicKey) IsEqual(otherPub *PublicKey) bool {

	if pub.CurveType != otherPub.CurveType {
		return false
	}

	if pub.Key.X.Cmp(otherPub.Key.X) != 0 || pub.Key.Y.Cmp(otherPub.Key.Y) != 0 {
		return false
	}

	return true
}

func (pub *PublicKey) IsEqualForTest(otherPub *PublicKey) bool {

	if pub.CurveType != otherPub.CurveType {
		return false
	}

	if pub.Key.X.Cmp(otherPub.Key.X) != 0 || pub.Key.Y.Cmp(otherPub.Key.Y) != 0 {
		return false
	}

	if !bytes.Equal(pub.Chaincode, otherPub.Chaincode) {
		return false
	}

	if !bytes.Equal(pub.RootFingerPrint, otherPub.RootFingerPrint) {
		return false
	}

	if pub.Index.Cmp(otherPub.Index) != 0 {
		return false
	}

	return true
}

func (pub *PublicKey) ToECDSA() *ecdsa.PublicKey {
	return pub.Key
}

func GetPublicKeyRootFingerprint(pub *ecdsa.PublicKey) ([]byte, error) {

	xLen := len(pub.X.Bytes())
	yLen := len(pub.Y.Bytes())

	rawBytes := make([]byte, xLen+yLen)
	copy(rawBytes, pub.X.Bytes())
	copy(rawBytes[xLen:], pub.Y.Bytes())

	hash, err := utils.SHA256RIPEMD160(rawBytes)
	if err != nil {
		return nil, err
	}

	if len(hash) < PUBLICKEY_FINGERPRINT_LEN {
		return nil, errors.New("Wrong private key fingerprint length")
	}

	return hash[:PUBLICKEY_FINGERPRINT_LEN], nil
}

func compressPublicKey(x *big.Int, y *big.Int) []byte {
	var key bytes.Buffer

	// Write header; 0x2 for even y value; 0x3 for odd
	key.WriteByte(byte(0x2) + byte(y.Bit(0)))

	// Write X coord; Pad the key so x is aligned with the LSB. Pad size is key length - header size (1) - xBytes size
	xBytes := x.Bytes()
	for i := 0; i < (PublicKeyCompressedLength - 1 - len(xBytes)); i++ {
		key.WriteByte(0x0)
	}
	key.Write(xBytes)

	return key.Bytes()
}

// As described at https://crypto.stackexchange.com/a/8916
func expandPublicKey(key []byte, curveParams *elliptic.CurveParams) (*big.Int, *big.Int) {
	Y := big.NewInt(0)
	X := big.NewInt(0)
	X.SetBytes(key[1:])

	// y^2 = x^3 + ax^2 + b
	// a = 0
	// => y^2 = x^3 + b
	ySquared := big.NewInt(0)
	ySquared.Exp(X, big.NewInt(3), nil)
	ySquared.Add(ySquared, curveParams.B)

	Y.ModSqrt(ySquared, curveParams.P)

	Ymod2 := big.NewInt(0)
	Ymod2.Mod(Y, big.NewInt(2))

	signY := uint64(key[0]) - 2
	if signY != Ymod2.Uint64() {
		Y.Sub(curveParams.P, Y)
	}

	return X, Y
}
