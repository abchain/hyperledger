package wallet

import (
	abcrypto "hyperledger.abchain.org/core/crypto"
)

type Wallet interface {

	// Create private key
	NewPrivKey(accountID string) (*abcrypto.PrivateKey, error)

	// Import private key
	ImportPrivKey(accountID string, privkey string) error

	// Import private key
	ImportPrivateKey(accountID string, privkey *abcrypto.PrivateKey) error

	// Load private key
	LoadPrivKey(accountID string) (*abcrypto.PrivateKey, error)

	// Remove private key
	RemovePrivKey(accountID string) error

	// Rename account id
	Rename(old string, new string) error

	// List all private keys
	ListAll() (map[string]*abcrypto.PrivateKey, error)

	// Read private keys from file
	Load() error

	// Write private keys to file
	Persist() error
}
