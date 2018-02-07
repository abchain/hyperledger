package crypto

import (
	"crypto/ecdsa"
	"io"
	"math/big"

	"errors"
	"hyperledger.abchain.org/protos"
)

var (
	// ErrInvalidPrivateKey is returned when a derived private key is invalid
	ErrInvalidPrivateKey = errors.New("Invalid private key")

	// ErrInvalidPublicKey is returned when a derived public key is invalid
	ErrInvalidPublicKey = errors.New("Invalid public key")
)

type Singer interface {
	Public() *PublicKey

	ChildPublic(index *big.Int) (*PublicKey, error)

	Sign(index *big.Int, rand io.Reader, hash []byte) (sig *Signature, err error)

	IsEqual(otherPriv *PrivateKey) bool

	ChildKey(index *big.Int) (*PrivateKey, error)

	Serialize() []byte

	Str() string

	ToECDSA() *ecdsa.PrivateKey

	PBMessage() *protos.PrivateKey
}

type Verifier interface {
	Verify(hash []byte, sig *Signature) bool

	IsEqual(otherPub *PublicKey) bool

	ChildKey(index *big.Int) (*PublicKey, error)

	Serialize() []byte

	Str() string

	ToECDSA() *ecdsa.PublicKey

	PBMessage() *protos.PublicKey
}
