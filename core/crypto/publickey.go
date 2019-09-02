package crypto

import (
	"fmt"
	"hyperledger.abchain.org/protos"
	"strings"
)

func DecodeCompactPublicKey(pk string) (Verifier, error) {
	pkparts := strings.Split(pk, ":")
	if len(pkparts) != 2 {
		return nil, fmt.Errorf("Invalid encoding for pk")
	}

	retpb := new(protos.PublicKey)
	retpb.Version = DefaultVersion

	switch pkparts[0] {
	case "EC":
		var ctype int32
		var key []byte

		if n, err := fmt.Sscanf(pkparts[1], "%d,%x", &ctype, &key); err != nil || n != 2 {
			return nil, fmt.Errorf("Decode ecdsa public key fail: %s [%d]", err, n)
		}

		if len(key) != 64 {
			return nil, fmt.Errorf("Decode ecdsa public key fail, invalid publickey length (%d)", len(key))
		}

		retpb.Pub = &protos.PublicKey_Ec{&protos.PublicKey_ECDSA{
			Curvetype: ctype,
			P:         &protos.ECPoint{X: key[:32], Y: key[32:]}}}
		ret := CryptoSchemes[SCHEME_ECDSA].NewVerifier()
		if err := ret.FromPBMessage(retpb); err != nil {
			return nil, err
		}
		return ret, nil
	default:
		return nil, fmt.Errorf("Unknown pk scheme: %s", pkparts[0])
	}
}
