package abchainTx

import (
	"hyperledger.abchain.org/core/crypto"
	pb "hyperledger.abchain.org/protos"
)

type AddrCredentials interface {
	Verify(addr Address) error
	CredCount() int
	GetCredPubkey(addr Address) *crypto.PublicKey
	ListCredPubkeys() []*crypto.PublicKey
}

type AddrCredentialBuilder interface {
	AddSignature(pub *crypto.PublicKey, sign *crypto.ECSignature)
	AddCc(ccname string, addr Address, pub *crypto.PublicKey) //is deprecated
	Update(msg *pb.TxCredential) error
}