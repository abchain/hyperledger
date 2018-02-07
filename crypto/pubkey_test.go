package crypto

import (
	"math/big"
	"testing"
)

func TestPublicKey_Serialize_SECP256K1(t *testing.T) {

	publicKeySerialize(t, SECP256K1, true)
}

func TestPublicKey_Serialize_ECP256_FIPS186(t *testing.T) {

	publicKeySerialize(t, ECP256_FIPS186, true)
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
		publicKeySerializeBenchmark(b, pub)
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
		publicKeySerializeBenchmark(b, pub)
	}
}

func publicKeySerialize(t *testing.T, curveType int32, log bool) {

	priv, err := NewPrivatekey(curveType)
	if err != nil {
		t.Fatal("Generate private key fail: ", err)
	}

	pub := priv.Public()
	raw := pub.Serialize()
	pub2, err := PublicKeyFromBytes(raw)
	if err != nil {
		t.Fatal("Import public key fail: ", err)
	}

	//t.Logf("Public key: %v", pub)
	//t.Logf("New Public key: %v", pub2)

	if !pub.IsEqualForTest(pub2) {
		t.Fatal("Public keys not equal")
	}

	if log {
		t.Logf("Dump public key: %v", pub.Str())
	}
}

func publicKeySerializeBenchmark(b *testing.B, pub *PublicKey) {

	raw := pub.Serialize()

	pub2, err := PublicKeyFromBytes(raw)
	if err != nil {
		b.Fatal("Import public key fail: ", err)
	}

	if !pub.IsEqualForTest(pub2) {
		b.Fatal("Public keys not equal")
	}
}
