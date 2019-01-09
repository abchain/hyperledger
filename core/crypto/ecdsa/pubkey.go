package ecdsa

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/core/crypto"
	"hyperledger.abchain.org/protos"
	"math/big"
)

type PublicKey struct {
	Version int32

	// Curve type
	CurveType int32

	// underlying public key
	Key *ecdsa.PublicKey

	*KeyDerivation
}

func (pub *PublicKey) PBMessage() proto.Message {

	// var X []byte
	// var Y []byte

	// // Compress algorithm only works for SECP256K1
	// if pub.CurveType == SECP256K1 {
	// 	X = compressPublicKey(pub.Key.X, pub.Key.Y)
	// 	Y = nil
	// } else {
	// 	X = pub.Key.X.Bytes()
	// 	Y = pub.Key.Y.Bytes()
	// }

	pubProto := new(protos.PublicKey)
	if pub.KeyDerivation != nil {
		pubProto.Kd = pub.ToKDMessage()
	}

	pubProto.Version = DefaultVersion

	pubProto.Pub = &protos.PublicKey_Ec{&protos.PublicKey_ECDSA{
		Curvetype: pub.CurveType,
		P:         &protos.ECPoint{pub.Key.X.Bytes(), pub.Key.Y.Bytes()},
	}}

	return pubProto
}

func (pub *PublicKey) resetPoint(p *protos.ECPoint) error {

	if p == nil {
		return fmt.Errorf("ECPoint is empty")
	}

	curve, err := GetEC(pub.CurveType)
	if err != nil {
		return err
	}

	// // Compress algorithm only works for SECP256K1
	// if pubProto.Curvetype == SECP256K1 {
	// 	X, Y = expandPublicKey(pubProto.GetP().X, curve.Params())
	// } else {
	// 	X.SetBytes(pubProto.GetP().X)
	// 	Y.SetBytes(pubProto.GetP().Y)
	// }

	pub.Key = &ecdsa.PublicKey{
		curve,
		big.NewInt(0).SetBytes(p.GetX()),
		big.NewInt(0).SetBytes(p.GetY()),
	}

	return nil
}

func (pub *PublicKey) FromPBMessage(msg proto.Message) error {
	var pubProto *protos.PublicKey_ECDSA

	if kmsg, ok := msg.(*protos.PublicKey); !ok {
		return fmt.Errorf("Not publickey message")
	} else if pubProto = kmsg.GetEc(); pubProto == nil {
		return fmt.Errorf("public key type [%v] is not ecdsa", kmsg.GetPub())
	} else if kmsg.GetVersion() != DefaultVersion {
		return fmt.Errorf("Unknown key version %d, (expect %d)", kmsg.GetVersion(), DefaultVersion)
	} else {

		pub.Version = kmsg.GetVersion()
		pub.CurveType = pubProto.GetCurvetype()

		if driv := kmsg.GetKd(); driv != nil {
			pub.KeyDerivation = new(KeyDerivation)
			pub.FromKDMessage(driv)
		}
	}

	return pub.resetPoint(pubProto.GetP())
}

// func (pub *PublicKey) Str() string {
// 	return base64.URLEncoding.EncodeToString(pub.Serialize())
// }

// func (pub *PublicKey) RootFingerPrintStr() string {

// 	return base64.URLEncoding.EncodeToString(pub.RootFingerPrint)
// }

func (pub *PublicKey) String() string {

	out := fmt.Sprintf("ecdsa_PublicKey{Version: %v, CurveType: %d", pub.Version, pub.CurveType)

	if pub.Key != nil {
		out = fmt.Sprintf("%s, X: %v, Y: %v", out, pub.Key.X, pub.Key.Y)
	}

	if pub.KeyDerivation != nil {
		out = fmt.Sprintf("%s, RootFingerPrint: %x, Index: %v, Chaincode: %x",
			out, pub.RootFingerPrint, pub.Index, pub.Chaincode)
	}

	return out + "}"
}

func (pub *PublicKey) child(index *big.Int) (*PublicKey, error) {

	if index.Int64() == 0 {
		return pub, nil
	}

	return getChildPublicKey(pub, index)
}

func (pub *PublicKey) Child(index *big.Int) (crypto.Hierarchical, error) {

	return pub.child(index)
}

func (pub *PublicKey) Verify(hash []byte, sig *protos.Signature) bool {

	ecsig := sig.GetEc()
	if ecsig == nil {
		return false
	}

	return ecdsa.Verify(pub.Key, hash, big.NewInt(0).SetBytes(ecsig.GetR()), big.NewInt(0).SetBytes(ecsig.GetS()))
}

func (pub *PublicKey) Recover(sig *protos.Signature) error {

	ecsig := sig.GetEc()
	if ecsig == nil {
		return fmt.Errorf("No ecdsa signature")
	} else if pk := ecsig.GetP(); pk == nil {
		return fmt.Errorf("recover data is not supported yet")
	} else {
		pub.KeyDerivation = nil
		pub.Version = DefaultVersion
		pub.CurveType = ecsig.GetCurvetype()

		return pub.resetPoint(pk)
	}
}

func (pub *PublicKey) IsEqual(p interface{}) bool {

	if otherPub, ok := p.(*PublicKey); !ok {
		return false
	} else if pub.CurveType != otherPub.CurveType {
		return false
	} else if pub.Key.X.Cmp(otherPub.Key.X) != 0 || pub.Key.Y.Cmp(otherPub.Key.Y) != 0 {
		return false
	}
	return true
}

//for testing and legacy purpose

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

const (
	PublicKeyCompressedLength = 33
)

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
