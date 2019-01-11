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

func TestCompatible(t *testing.T) {

	priv := new(PrivateKey)
	priv.Version = DefaultVersion
	priv.CurveType = SECP256K1
	priv.Key = genPrivkey(priv.CurveType, "110871560887871433772251872457573575186844291339401824426745393239615136047255")
	priv.KeyDerivation = new(KeyDerivation)
	priv.Index = big.NewInt(0)
	priv.Chaincode = []byte{0x5b, 0x3c, 0x57, 0x84, 0xd2, 0x63, 0x68, 0x60, 0x75, 0xa3, 0x4f, 0x0f, 0x4f, 0x1a, 0xa2, 0x39, 0xf3, 0x78, 0x40, 0x5a, 0x1e, 0x04, 0xf9, 0x72, 0x33, 0x55, 0x2b, 0xbf, 0x75, 0xf4, 0x25, 0x23}

	childPriv, err := priv.child(big.NewInt(0))
	if err != nil {
		t.Fatalf("Child fail", err)
	}

	expectedCD, _ := big.NewInt(0).SetString("70912661286874047882889983861262634846432069668246130355602184839087556488003", 10)
	expectedCC := []byte{0x77, 0xf0, 0x40, 0x72, 0xb1, 0x08, 0x4f, 0xdc, 0xd8, 0x5e, 0xfd, 0xbf, 0xd7, 0x34, 0xd9, 0xe6, 0x25, 0x62, 0x18, 0xa0, 0xb4, 0x66, 0x82, 0x07, 0xe9, 0xcc, 0x96, 0x28, 0x4f, 0x07, 0x1a, 0x8a}

	if childPriv.Key.D.Cmp(expectedCD) != 0 {
		t.Fatalf("Unexpected priv (%s, expect %s)", childPriv.Key.D, expectedCD)
	}

	if bytes.Compare(expectedCC, childPriv.Chaincode) != 0 {
		t.Fatalf("Unexpected cc (%x, expect %x)", childPriv.Chaincode, expectedCC)
	}

	childPriv2, err := childPriv.child(big.NewInt(1))
	if err != nil {
		t.Fatalf("Child 2 fail", err)
	}

	expectedCD, _ = big.NewInt(0).SetString("26450297871391175399259363899116702850686174311172247471930977743122743540079", 10)
	expectedCC = []byte{0xa1, 0x21, 0x19, 0x82, 0x15, 0xfd, 0x38, 0x6f, 0x46, 0x28, 0x5d, 0xd0, 0x0b, 0xf7, 0x2e, 0x09, 0x10, 0xea, 0xde, 0x7b, 0xb7, 0xe9, 0x7e, 0xc5, 0xdd, 0x3a, 0x16, 0x81, 0xe6, 0x2b, 0xea, 0x6f}

	if childPriv2.Key.D.Cmp(expectedCD) != 0 {
		t.Fatalf("Unexpected priv 2 (%s, expect %s)", childPriv2.Key.D, expectedCD)
	}

	if bytes.Compare(expectedCC, childPriv2.Chaincode) != 0 {
		t.Fatalf("Unexpected cc (%x, expect %x)", childPriv2.Chaincode, expectedCC)
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
