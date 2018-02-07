package abchainTx

import (
	"math/big"
	"testing"

	abcrypto "hyperledger.abchain.org/crypto"
)

func GetPublicKey(t *testing.T) *abcrypto.PublicKey {
	priv, err := abcrypto.NewPrivatekey(abcrypto.DefaultCurveType)
	if err != nil {
		t.Fatal("Generate private key fail: ", err)
	}

	pub, err := priv.ChildPublic(big.NewInt(0))
	if err != nil {
		t.Fatal("Get root public key fail: ", err)
	}

	return pub
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
