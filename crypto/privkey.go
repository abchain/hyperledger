package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/protos"
	"io"
	"math/big"
)

const (
	PRIVATEKEY_VERSION = 1
)

type PrivateKey struct {
	Version int32

	// Curve type
	CurveType int32

	// underlying private key
	Key *ecdsa.PrivateKey

	// underlying root private key fingerprint
	RootFingerPrint []byte

	// index of child private key
	Index *big.Int

	// chaincode
	Chaincode []byte
}

var (
	DefaultPrivateKeyVersion int32 = PRIVATEKEY_VERSION
)

func NewPrivatekey(curveType int32) (*PrivateKey, error) {

	curve, err := GetEC(curveType)
	if err != nil {
		return nil, err
	}

	// Generate underlying private key
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	// Validate underlying private key
	if err = validatePrivateKey(key.D.Bytes(), curve.Params()); err != nil {
		return nil, err
	}

	// Generate chaincode
	chaincodeKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	// Generate Root Fingerprint
	rootFingerprint, err := GetPublicKeyRootFingerprint(&key.PublicKey)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{
		DefaultPrivateKeyVersion,
		curveType,
		key,
		rootFingerprint,
		big.NewInt(0),
		chaincodeKey.D.Bytes()}, nil
}

func PrivatekeyFromString(privkeyStr string) (*PrivateKey, error) {

	raw, err := base64.URLEncoding.DecodeString(privkeyStr)
	if err != nil {
		return nil, err
	}

	return PrivatekeyFromBytes(raw)
}

func PrivatekeyFromBytes(raw []byte) (*PrivateKey, error) {

	privProto := &protos.PrivateKey{}
	err := proto.Unmarshal(raw, privProto)
	if err != nil {
		return nil, err
	}

	return PrivateKeyFromPBMessage(privProto)
}

func PrivateKeyFromPBMessage(privProto *protos.PrivateKey) (*PrivateKey, error) {

	if privProto.Version != PRIVATEKEY_VERSION {
		return nil, errors.New("Unknown private key version")
	}

	curve, err := GetEC(privProto.Curvetype)
	if err != nil {
		return nil, err
	}

	x, y := curve.ScalarBaseMult(privProto.D)

	ecdsaKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(privProto.D),
	}

	return &PrivateKey{
		DefaultPrivateKeyVersion,
		privProto.Curvetype,
		ecdsaKey,
		privProto.RootFingerprint,
		big.NewInt(0).SetBytes(privProto.Index),
		privProto.Chaincode}, nil
}

func (priv *PrivateKey) PBMessage() *protos.PrivateKey {

	privProto := &protos.PrivateKey{
		priv.Version,
		priv.CurveType,
		priv.RootFingerPrint,
		priv.Index.Bytes(),
		priv.Chaincode,
		priv.Key.D.Bytes(),
	}

	return privProto
}

func (priv *PrivateKey) Serialize() []byte {

	privProto := priv.PBMessage()

	raw, err := proto.Marshal(privProto)
	if err != nil {
		return nil
	}

	return raw
}

func (priv *PrivateKey) Str() string {
	return base64.URLEncoding.EncodeToString(priv.Serialize())
}

func (priv *PrivateKey) String() string {

	return fmt.Sprintf("&{Version: %v, CurveType: %d, Key.D: %v, "+
		"RootFingerPrint: %v, Index: %v, Chaincode: %v}",
		priv.Version, priv.CurveType, priv.Key.D.Bytes(),
		priv.RootFingerPrint, priv.Index, priv.Chaincode)
}

func (priv *PrivateKey) Public() *PublicKey {

	pubkey := &PublicKey{
		DefaultPublicKeyVersion,
		priv.CurveType,
		&priv.Key.PublicKey,
		priv.RootFingerPrint,
		priv.Index,
		priv.Chaincode}

	return pubkey
}

func (priv *PrivateKey) ChildKey(index *big.Int) (*PrivateKey, error) {

	if index.Int64() == 0 {
		return priv, nil
	}

	return getChildPrivateKey(priv, index)
}

func (priv *PrivateKey) ChildPublic(index *big.Int) (*PublicKey, error) {

	childPrivkey, err := priv.ChildKey(index)
	if err != nil {
		return nil, err
	}

	return childPrivkey.Public(), nil
}

func (priv *PrivateKey) Sign(index *big.Int, rand io.Reader, hash []byte) (sig *Signature, err error) {

	childKey, err := priv.ChildKey(index)
	if err != nil {
		return nil, err
	}

	r, s, err := ecdsa.Sign(rand, childKey.Key, hash)

	return &Signature{r, s}, err
}

func (priv *PrivateKey) SignwithThis(rand io.Reader, hash []byte) (sig *Signature, err error) {

	r, s, err := ecdsa.Sign(rand, priv.Key, hash)

	return &Signature{r, s}, err
}

func (priv *PrivateKey) IsEqual(otherPriv *PrivateKey) bool {

	if priv.Key.D.Cmp(otherPriv.Key.D) != 0 {
		return false
	}

	return true
}

func (priv *PrivateKey) IsEqualForTest(otherPriv *PrivateKey) bool {

	if priv.Key.D.Cmp(otherPriv.Key.D) != 0 {
		return false
	}

	if !bytes.Equal(priv.Chaincode, otherPriv.Chaincode) {
		return false
	}

	if !bytes.Equal(priv.RootFingerPrint, otherPriv.RootFingerPrint) {
		return false
	}

	if priv.Index.Cmp(otherPriv.Index) != 0 {
		return false
	}

	return true
}

func (priv *PrivateKey) ToECDSA() *ecdsa.PrivateKey {
	return priv.Key
}

func GetPrivateKeyRootFingerprint(priv *ecdsa.PrivateKey) ([]byte, error) {

	return GetPublicKeyRootFingerprint(&priv.PublicKey)
}

func validatePrivateKey(key []byte, curveParams *elliptic.CurveParams) error {
	if fmt.Sprintf("%x", key) == "0000000000000000000000000000000000000000000000000000000000000000" ||
		bytes.Compare(key, curveParams.N.Bytes()) >= 0 {
		return ErrInvalidPrivateKey
	}

	return nil
}
