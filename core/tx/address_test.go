package abchainTx

import (
	"testing"

	abcrypto "hyperledger.abchain.org/core/crypto"
	"hyperledger.abchain.org/core/crypto/ecdsa"
)

func GetPublicKey(t *testing.T) abcrypto.Verifier {
	priv, err := ecdsa.NewPrivatekey(ecdsa.DefaultCurveType)
	if err != nil {
		t.Fatal("Generate private key fail: ", err)
	}

	return priv.Public()
}

func TestAddress_Serialize(t *testing.T) {

	pub := GetPublicKey(t)

	addr, err := NewAddress(pub)
	if err != nil {
		t.Fatal("Get address fail")
	}

	addrStr := addr.ToString()
	t.Logf("Generate Address: %v", addrStr)

	addr2, err := NewAddressFromString(addrStr)
	if err != nil {
		t.Fatal("Import address fail: ", err)
	}

	t.Logf("Import Address: %v", addr2.ToString())

	if !addr.IsEqual(addr2) {
		t.Fatal("Addresses not equal: ")
	}
}
