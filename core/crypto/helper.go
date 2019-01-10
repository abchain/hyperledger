package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	proto "github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/protos"
	"io"
	"math/big"
)

var DefaultRandSrc io.Reader
var CryptoSchemes = map[string]Factory{}

func init() {
	DefaultRandSrc = rand.Reader
}

func PublicKeyFromBytes(raw []byte) (Verifier, error) {

	pubProto := &protos.PublicKey{}
	err := proto.Unmarshal(raw, pubProto)
	if err != nil {
		return nil, err
	}

	return PublicKeyFromPBMessage(pubProto)
}

func PublicKeyToBytes(pk Verifier) ([]byte, error) {

	return proto.Marshal(pk.PBMessage())
}

func PrivatekeyFromString(privkeyStr string) (Signer, error) {

	raw, err := base64.URLEncoding.DecodeString(privkeyStr)
	if err != nil {
		return nil, err
	}

	return PrivatekeyFromBytes(raw)
}

func PrivatekeyToString(priv Signer) (string, error) {

	raw, err := proto.Marshal(priv.PBMessage())
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(raw), nil
}

func PrivatekeyFromBytes(raw []byte) (Signer, error) {

	privProto := &protos.PrivateKey{}
	err := proto.Unmarshal(raw, privProto)
	if err != nil {
		return nil, err
	}

	return PrivateKeyFromPBMessage(privProto)
}

func PublicKeyFromPBMessage(pubProto *protos.PublicKey) (Verifier, error) {

	switch pubProto.Pub.(type) {
	case *protos.PublicKey_Ec:
		fac, ok := CryptoSchemes[SCHEME_ECDSA]
		if !ok {
			return nil, errors.New("ECDSA scheme is not available")
		}

		v := fac.NewVerifier()
		if err := v.FromPBMessage(pubProto); err != nil {
			return nil, err
		}

		return v, nil
	default:
		return nil, errors.New("Unknown public key type")
	}

}

func PrivateKeyFromPBMessage(privProto *protos.PrivateKey) (Signer, error) {

	switch privProto.Priv.(type) {
	case *protos.PrivateKey_Ec:
		fac, ok := CryptoSchemes[SCHEME_ECDSA]
		if !ok {
			return nil, errors.New("ECDSA scheme is not available")
		}

		v := fac.NewSigner()
		if err := v.FromPBMessage(privProto); err != nil {
			return nil, err
		}

		return v, nil
	default:
		return nil, errors.New("Unknown private key type")
	}

}

func PublicKeyFromSignature(sig *protos.Signature) (Verifier, error) {

	switch sig.GetData().(type) {
	case *protos.Signature_Ec:
		fac, ok := CryptoSchemes[SCHEME_ECDSA]
		if !ok {
			return nil, errors.New("ECDSA scheme is not available")
		}

		v := fac.NewVerifier()
		if err := v.Recover(sig); err != nil {
			return nil, err
		}

		return v, nil
	default:
		return nil, errors.New("Unknown signature type")
	}

}

func PrivateKeySign(s Signer, index *big.Int, hash []byte) (*protos.Signature, error) {
	if chs, err := GetChildPrivateKey(s, index); err != nil {
		return nil, err
	} else {
		return chs.Sign(hash)
	}
}

func GetChildPrivateKey(s Signer, index *big.Int) (Signer, error) {
	if ch, err := s.Child(index); err != nil {
		return nil, err
	} else if chs, ok := ch.(Signer); !ok {
		return nil, errors.New("Assigned a child not qualified as an signer")
	} else {
		return chs, nil
	}
}

func GetChildPublicKey(pk Verifier, index *big.Int) (Verifier, error) {
	if ch, err := pk.Child(index); err != nil {
		return nil, err
	} else if chpk, ok := ch.(Verifier); !ok {
		return nil, errors.New("Assigned a child not qualified as an verifier")
	} else {
		return chpk, nil
	}
}
