package crypto

import (
	"crypto/ecdsa"
	"errors"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/protos"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

type ECSignature struct {
	Signature
}

func SignatureFromBytes(raw []byte) (*Signature, error) {

	sigProto := &protos.Signature{}
	err := proto.Unmarshal(raw, sigProto)
	if err != nil {
		return nil, err
	}

	return SignatureFromPBMessage(sigProto)
}

func SignatureFromPBMessage(sigProto *protos.Signature) (*Signature, error) {

	if sigProto == nil {
		return nil, errors.New("SignatureFromPBMessage: input null pointer")
	}

	sig := &Signature{
		big.NewInt(0).SetBytes(sigProto.GetP().X),
		big.NewInt(0).SetBytes(sigProto.GetP().Y),
	}

	return sig, nil
}

func ECSignatureFromPB(sig *protos.ECSignature) (*ECSignature, error) {

	if sig == nil {
		return nil, errors.New("ECSignatureFromPB: input null pointer")
	}

	ecs := &ECSignature{
		Signature{
			big.NewInt(0).SetBytes(sig.R),
			big.NewInt(0).SetBytes(sig.S),
		},
	}

	return ecs, nil
}

func (sig *Signature) PBMessage() *protos.Signature {

	sigProto := &protos.Signature{
		&protos.ECPoint{
			sig.R.Bytes(),
			sig.S.Bytes(),
		},
	}

	return sigProto
}

func (sig *Signature) Serialize() []byte {

	sigProto := sig.PBMessage()

	raw, err := proto.Marshal(sigProto)
	if err != nil {
		return nil
	}

	return raw
}

// Verify calls ecdsa.Verify to verify the signature of hash using the public
// key.  It returns true if the signature is valid, false otherwise.
func (sig *Signature) Verify(hash []byte, pubKey *PublicKey) bool {
	return ecdsa.Verify(pubKey.ToECDSA(), hash, sig.R, sig.S)
}

// IsEqual compares this Signature instance to the one passed, returning true
// if both Signatures are equivalent. A signature is equivalent to another, if
// they both have the same scalar value for R and S.
func (sig *Signature) IsEqual(otherSig *Signature) bool {
	return sig.R.Cmp(otherSig.R) == 0 &&
		sig.S.Cmp(otherSig.S) == 0
}

func (s *ECSignature) ToPB() *protos.ECSignature {

	return &protos.ECSignature{
		s.R.Bytes(),
		s.S.Bytes(),
		0,
	}
}

func (sig *ECSignature) Serialize() []byte {

	sigProto := sig.ToPB()

	raw, err := proto.Marshal(sigProto)
	if err != nil {
		return nil
	}

	return raw
}
