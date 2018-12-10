package msp

import (
	"fmt"

	"hyperledger.abchain.org/client/hyfabric/utils"

	//mspmgmt "github.com/hyperledger/fabric/msp/mgmt"
	cb "github.com/hyperledger/fabric/protos/common"
)

// LocalSigner is a temporary stub interface which will be implemented by the local MSP
type LocalSigner interface {
	SignatureHeaderMaker
	Signer
}

// Signer signs messages
type Signer interface {
	// Sign a message and return the signature over the digest, or error on failure
	Sign(message []byte) ([]byte, error)
}

// IdentitySerializer serializes identities
type IdentitySerializer interface {
	// Serialize converts an identity to bytes
	Serialize() ([]byte, error)
}

// SignatureHeaderMaker creates a new SignatureHeader
type SignatureHeaderMaker interface {
	// NewSignatureHeader creates a SignatureHeader with the correct signing identity and a valid nonce
	NewSignatureHeader() (*cb.SignatureHeader, error)
}

type mspSigner struct {
}

// NewSigner returns a new instance of the msp-based LocalSigner.
// It assumes that the local msp has been already initialized.
// Look at mspmgmt.LoadLocalMsp for further information.
// func NewSigner() crypto.LocalSigner {
// 	return &mspSigner{}
// }
func NewSigner() LocalSigner {
	return &mspSigner{}
}

// NewSignatureHeader creates a SignatureHeader with the correct signing identity and a valid nonce
func (s *mspSigner) NewSignatureHeader() (*cb.SignatureHeader, error) {
	logger.Debug("msp NewSignatureHeader===========")
	signer, err := GetSigningIdentity()
	if err != nil {
		return nil, fmt.Errorf("Failed getting MSP-based signer [%s]", err)
	}

	//实现：github.com\hyperledger\fabric\msp\identities.go
	creatorIdentityRaw, err := signer.Serialize()

	if err != nil {
		return nil, fmt.Errorf("Failed serializing creator public identity [%s]", err)
	}

	nonce, err := utils.GetRandomNonce()
	if err != nil {
		return nil, fmt.Errorf("Failed creating nonce [%s]", err)
	}

	sh := &cb.SignatureHeader{}
	sh.Creator = creatorIdentityRaw
	sh.Nonce = nonce

	return sh, nil
}

// Sign a message which should embed a signature header created by NewSignatureHeader
func (s *mspSigner) Sign(message []byte) ([]byte, error) {
	signer, err := GetSigningIdentity()
	if err != nil {
		return nil, fmt.Errorf("Failed getting MSP-based signer [%s]", err)
	}

	signature, err := signer.Sign(message)
	if err != nil {
		return nil, fmt.Errorf("Failed generating signature [%s]", err)
	}

	return signature, nil
}
