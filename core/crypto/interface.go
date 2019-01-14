package crypto

import (
	proto "github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/protos"
	"math/big"
)

const (
	SCHEME_ECDSA              = "ecdsa"
	PUBLICKEY_FINGERPRINT_LEN = 8
)

var DefaultVersion int32

type Factory interface {
	NewSigner() Signer
	NewVerifier() Verifier
}

type Signer interface {
	base
	Hierarchical
	Public() Verifier
	Sign([]byte) (*protos.Signature, error)
}

type Verifier interface {
	base
	Hierarchical
	Verify([]byte, *protos.Signature) bool
	//the Verifier recovered from signature is not hierarchical
	Recover(*protos.Signature) error
	Digest() []byte
}

type base interface {
	IsEqual(interface{}) bool //not compare the hierarchical data
	String() string
	PBMessage() proto.Message
	FromPBMessage(proto.Message) error
}

type Hierarchical interface {
	GetRootFingerPrint() []byte
	GetIndex() *big.Int
	Child(*big.Int) (Hierarchical, error)
}
