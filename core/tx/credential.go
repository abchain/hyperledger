package abchainTx

import (
	"hyperledger.abchain.org/core/crypto"
	pb "hyperledger.abchain.org/protos"
)

type AddrCredentials interface {
	Verify(addr Address) error
	CredCount() int
	// GetCredPubkey(addr Address) *crypto.PublicKey
	// ListCredPubkeys() []*crypto.PublicKey
	GetCredPubkey(addr Address) crypto.Verifier
	ListCredPubkeys() []crypto.Verifier
}

type AddrCredentialBuilder interface {
	// AddSignature(pub *crypto.PublicKey, sign *crypto.ECSignature)
	// AddCc(ccname string, addr Address, pub *crypto.PublicKey) //is deprecated
	AddSignature(sign *pb.Signature)
	AddCc(ccname string, addr Address, pub crypto.Verifier) //is deprecated
	Update(msg *pb.TxCredential) error
}
