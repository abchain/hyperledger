package crypto

import (
	"crypto/rand"
	"errors"
	proto "github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/protos"
	"io"
	"math/big"
)

var DefaultRandSrc io.Reader
var DefaultCryptoScheme Factory
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

func PrivateKeySign(s Signer, index *big.Int, hash []byte) (*protos.Signature, error) {
	if ch, err := s.Child(index); err != nil {
		return nil, err
	} else if chs, ok := ch.(Signer); !ok {
		return nil, errors.New("Assigned a child not qualified as an signer")
	} else {
		return chs.Sign(hash)
	}
}
