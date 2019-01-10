package ecdsa

import (
	"bytes"
	"crypto/ecdsa"

	"fmt"
	proto "github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/core/crypto"
	"hyperledger.abchain.org/protos"
	"io"
	"math/big"
)

type PrivateKey struct {
	Version int32

	// Curve type
	CurveType int32

	// underlying private key
	Key *ecdsa.PrivateKey

	*KeyDerivation
}

func NewDefaultPrivatekey() (*PrivateKey, error) {
	return NewPrivatekey(DefaultCurveType)
}

func NewPrivatekey(curveType int32) (*PrivateKey, error) {

	curve, err := GetEC(curveType)
	if err != nil {
		return nil, err
	}

	// Generate underlying private key
	key, err := ecdsa.GenerateKey(curve, crypto.DefaultRandSrc)
	if err != nil {
		return nil, err
	}
	// Validate underlying private key
	if err = validatePrivateKey(key.D.Bytes(), curve.Params()); err != nil {
		return nil, err
	}

	// chaincode is 256bit
	cc := make([]byte, 32)

	if _, err := io.ReadFull(crypto.DefaultRandSrc, cc); err != nil {
		return nil, err
	}

	return &PrivateKey{
		DefaultVersion,
		curveType,
		key,
		&KeyDerivation{
			nil,
			big.NewInt(0),
			cc,
		}}, nil
}

// func PrivatekeyFromString(privkeyStr string) (*PrivateKey, error) {

// 	raw, err := base64.URLEncoding.DecodeString(privkeyStr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return PrivatekeyFromBytes(raw)
// }

// func PrivatekeyFromBytes(raw []byte) (*PrivateKey, error) {

// 	privProto := &protos.PrivateKey{}
// 	err := proto.Unmarshal(raw, privProto)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return PrivateKeyFromPBMessage(privProto)
// }

// func PrivateKeyFromPBMessage(privProto *protos.PrivateKey) (*PrivateKey, error) {

// 	if privProto.Version != PRIVATEKEY_VERSION {
// 		return nil, fmt.Errorf("Unknown private key version")
// 	}

// 	curve, err := GetEC(privProto.Curvetype)
// 	if err != nil {
// 		return nil, err
// 	}

// 	x, y := curve.ScalarBaseMult(privProto.D)

// 	ecdsaKey := &ecdsa.PrivateKey{
// 		PublicKey: ecdsa.PublicKey{
// 			Curve: curve,
// 			X:     x,
// 			Y:     y,
// 		},
// 		D: new(big.Int).SetBytes(privProto.D),
// 	}

// 	return &PrivateKey{
// 		DefaultPrivateKeyVersion,
// 		privProto.Curvetype,
// 		ecdsaKey,
// 		privProto.RootFingerprint,
// 		big.NewInt(0).SetBytes(privProto.Index),
// 		privProto.Chaincode}, nil
// }

func (priv *PrivateKey) PBMessage() proto.Message {

	privProto := new(protos.PrivateKey)

	if priv.KeyDerivation != nil {
		privProto.Kd = priv.ToKDMessage()
	}

	privProto.Version = DefaultVersion
	privProto.Priv = &protos.PrivateKey_Ec{&protos.PrivateKey_ECDSA{
		Curvetype: priv.CurveType,
		D:         priv.Key.D.Bytes(),
	}}

	return privProto
}

func (priv *PrivateKey) FromPBMessage(msg proto.Message) error {

	var privProto *protos.PrivateKey_ECDSA

	if kmsg, ok := msg.(*protos.PrivateKey); !ok {
		return fmt.Errorf("Not private message")
	} else if privProto = kmsg.GetEc(); privProto == nil {
		return fmt.Errorf("priv key type [%v] is not ecdsa", kmsg.GetPriv())
	} else if kmsg.GetVersion() != DefaultVersion {
		return fmt.Errorf("Unknown key version %d, (expect %d)", kmsg.GetVersion(), DefaultVersion)
	} else {

		priv.Version = kmsg.GetVersion()
		priv.CurveType = privProto.GetCurvetype()

		if driv := kmsg.GetKd(); driv != nil {
			priv.KeyDerivation = new(KeyDerivation)
			priv.FromKDMessage(driv)
		}
	}

	curve, err := GetEC(priv.CurveType)
	if err != nil {
		return err
	}

	x, y := curve.ScalarBaseMult(privProto.D)

	priv.Key = &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(privProto.D),
	}

	return nil
}

// func (priv *PrivateKey) Serialize() []byte {

// 	privProto := priv.PBMessage()

// 	raw, err := proto.Marshal(privProto)
// 	if err != nil {
// 		return nil
// 	}

// 	return raw
// }

// func (priv *PrivateKey) Str() string {
// 	return base64.URLEncoding.EncodeToString(priv.Serialize())
// }

func (priv *PrivateKey) String() string {

	out := fmt.Sprintf("ecdsa_PrivateKey{Version: %v, CurveType: %d", priv.Version, priv.CurveType)

	if priv.Key != nil {
		out = fmt.Sprintf("%s, D: %v", out, priv.Key.D)
	}

	if priv.KeyDerivation != nil {
		out = fmt.Sprintf("%s, RootFingerPrint: %x, Index: %v, Chaincode: %x",
			out, priv.RootFingerPrint, priv.Index, priv.Chaincode)
	}

	return out + "}"

}

func (priv *PrivateKey) public() *PublicKey {

	return &PublicKey{
		priv.Version,
		priv.CurveType,
		&priv.Key.PublicKey,
		priv.KeyDerivation,
	}

}

func (priv *PrivateKey) Public() crypto.Verifier {

	return priv.public()

}

func (priv *PrivateKey) child(index *big.Int) (*PrivateKey, error) {

	return getChildPrivateKey(priv, index)
}

func (priv *PrivateKey) Child(index *big.Int) (crypto.Hierarchical, error) {

	return priv.child(index)
}

func (priv *PrivateKey) IsEqual(p interface{}) bool {

	if otherPriv, ok := p.(*PrivateKey); !ok {
		return false
	} else if priv.CurveType != otherPriv.CurveType {
		return false
	} else if priv.Key.D.Cmp(otherPriv.Key.D) != 0 {
		return false
	}

	return true
}

func (priv *PrivateKey) Sign(hash []byte) (sig *protos.Signature, err error) {

	r, s, err := ecdsa.Sign(crypto.DefaultRandSrc, priv.Key, hash)

	if err != nil {
		return nil, err
	}

	//TODO: using v is still not available, we put pk
	sigpk := &protos.Signature_ECDSA_P{
		&protos.ECPoint{
			priv.Key.PublicKey.X.Bytes(),
			priv.Key.PublicKey.Y.Bytes(),
		},
	}

	ecsign := &protos.Signature_ECDSA{
		Curvetype: priv.CurveType,
		R:         r.Bytes(),
		S:         s.Bytes(),
		Pub:       sigpk,
	}

	var kdinf *protos.KeyDerived
	if priv.GetRootFingerPrint() != nil {
		kdinf = priv.ToKDMessage()
	}

	return &protos.Signature{&protos.Signature_Ec{ecsign}, kdinf}, nil
}

//deprecate APIs
func (priv *PrivateKey) ToECDSA() *ecdsa.PrivateKey {
	return priv.Key
}

//for testing and legacy purpose

func (priv *PrivateKey) ChildPublic(index *big.Int) (*PublicKey, error) {

	childPrivkey, err := priv.child(index)
	if err != nil {
		return nil, err
	}

	return childPrivkey.public(), nil
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
