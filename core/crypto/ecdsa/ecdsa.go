package ecdsa

import (
	"crypto/elliptic"
	"fmt"
	"hyperledger.abchain.org/core/crypto"
)

const (
	SECP256K1 int32 = iota + 1
	ECP256_FIPS186
)

type cryptoFactory struct{}

func (cryptoFactory) NewSigner() crypto.Signer {
	return new(PrivateKey)
}

func (cryptoFactory) NewVerifier() crypto.Verifier {
	return new(PublicKey)
}

var (
	DefaultCurveType = ECP256_FIPS186
	DefaultVersion   int32
	DefaultFactory   = cryptoFactory{}
)

func GetEC(curveType int32) (elliptic.Curve, error) {
	switch curveType {
	case ECP256_FIPS186:
		return elliptic.P256(), nil
	case SECP256K1:
		return Secp256k1(), nil
	default:
		return nil, fmt.Errorf("%d is not a valid curve defination", curveType)
	}
}

func init() {
	crypto.CryptoSchemes[crypto.SCHEME_ECDSA] = DefaultFactory
}
