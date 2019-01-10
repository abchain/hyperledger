package ecdsa

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"
)

func TestBIP32_SECP256K1(t *testing.T) {

	var err error

	// Generate private key
	priv, err := NewPrivatekey(DefaultCurveType)
	if err != nil {
		t.Fatal("Generate private key fail: ", err)
	}

	err = verifyBIP32(priv, big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}

	err = verifyBIP32(priv, big.NewInt(10000))
	if err != nil {
		t.Fatal(err)
	}

	err = verifyBIP32(priv, big.NewInt(10000000000))
	if err != nil {
		t.Fatal(err)
	}

	err = verifyBIP32(priv, big.NewInt(99999999999999999))
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkBIP32_Verify_Case01(b *testing.B) {

	var err error

	// Generate private key
	priv, err := NewPrivatekey(DefaultCurveType)
	if err != nil {
		b.Fatal("Generate private key fail: ", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		index := big.NewInt(int64(i))

		err = verifyBIP32(priv, index)
		if err != nil {
			b.Log("[%i] %v", index, err)
		}
	}
}

func BenchmarkBIP32_Verify_Case02(b *testing.B) {

	for i := 0; i < b.N; i++ {

		// Generate private key
		priv, err := NewPrivatekey(DefaultCurveType)
		if err != nil {
			b.Fatal("Generate private key fail: ", err)
		}

		index := big.NewInt(int64(i))

		err = verifyBIP32(priv, index)
		if err != nil {
			b.Log("[%i] %v", index, err)
		}
	}
}

func verifyBIP32(priv *PrivateKey, index *big.Int) error {

	rootpubkey := priv.public()

	childPrivkey, err := priv.child(index)
	if err != nil {
		return fmt.Errorf("Get child private key fail: ", err)
	}

	childPubkey1 := childPrivkey.Public()

	childPubkey2, err := rootpubkey.child(index)
	if err != nil {
		return fmt.Errorf("Get child public key1 fail: ", err)
	}

	if !childPubkey1.IsEqual(childPubkey2) {
		return fmt.Errorf("Child public keys not equal")
	}

	if bytes.Compare(childPubkey1.GetRootFingerPrint(), childPubkey2.GetRootFingerPrint()) != 0 {
		return fmt.Errorf("Child key's root fingerprint not equal:", childPubkey1, childPubkey2)
	}

	if bytes.Compare(rootpubkey.Digest()[:PUBLICKEY_FINGERPRINT_LEN], childPubkey1.GetRootFingerPrint()) != 0 {
		return fmt.Errorf("Child key's root fingerprint not equal to root:", rootpubkey.Digest(), childPubkey2)
	}

	if index.Cmp(big.NewInt(0)) != 0 && childPrivkey.IsEqual(priv) {
		return fmt.Errorf("Child private key is equal to root private key")
	}

	if index.Cmp(big.NewInt(0)) != 0 && childPubkey1.IsEqual(rootpubkey) {
		return fmt.Errorf("Child publick key is equal to root publick key")
	}

	return nil
}
