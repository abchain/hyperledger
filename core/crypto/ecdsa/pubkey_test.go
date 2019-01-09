package ecdsa

import (
	"math/big"
	"testing"
)

func TestPublicKey_Serialize_SECP256K1(t *testing.T) {

	priv, err := NewPrivatekey(SECP256K1)
	if err != nil {
		t.Fatal("Generate private key fail: ", err)
	}

	publicKeySerialize(t, priv.public())

	t.Logf("Dump public key: %s", priv.public())
}

func TestPublicKey_Serialize_ECP256_FIPS186(t *testing.T) {

	priv, err := NewPrivatekey(ECP256_FIPS186)
	if err != nil {
		t.Fatal("Generate private key fail: ", err)
	}

	publicKeySerialize(t, priv.public())

	t.Logf("Dump public key: %s", priv.public())

}

func BenchmarkPublicKey_Serialize_SECP256K1(b *testing.B) {

	priv, err := NewPrivatekey(SECP256K1)
	if err != nil {
		b.Fatal("Generate private key fail: ", err)
	}

	pub, err := priv.ChildPublic(big.NewInt(0))
	if err != nil {
		b.Fatal("Get root public key fail: ", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		publicKeySerialize(b, pub)
	}
}

func BenchmarkPublicKey_Serialize_ECP256_FIPS186(b *testing.B) {

	priv, err := NewPrivatekey(ECP256_FIPS186)
	if err != nil {
		b.Fatal("Generate private key fail: ", err)
	}

	pub, err := priv.ChildPublic(big.NewInt(0))
	if err != nil {
		b.Fatal("Get root public key fail: ", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		publicKeySerialize(b, pub)
	}
}

func publicKeySerialize(b testing.TB, pub *PublicKey) {

	pubpb := pub.PBMessage()

	pub2 := new(PublicKey)
	if err := pub2.FromPBMessage(pubpb); err != nil {
		b.Fatal("Import public key fail: ", err)
	}

	if !pub.IsEqualForTest(pub2) {
		b.Fatal("Public keys not equal", pub, pub2)
	}
}
