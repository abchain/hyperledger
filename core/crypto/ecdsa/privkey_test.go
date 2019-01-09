package ecdsa

import (
	"crypto/rand"
	"hyperledger.abchain.org/core/crypto"
	"math/big"
	"testing"
)

func TestSigning(t *testing.T) {

	// Generate private key

	priv, err := NewPrivatekey(DefaultCurveType)
	if err != nil {
		t.Fatal("Generate private key fail: ", err)
	}

	// Generate message data

	rb := make([]byte, 32)
	_, err = rand.Read(rb)
	if err != nil {
		t.Fatal("Genereate rand bytes fail: ", err)
	}

	// Root Private Key Signing

	sig, err := priv.Sign(rb)
	if err != nil {
		t.Fatal("Sign raw data fail: ", err)
	}

	pub, err := priv.ChildPublic(big.NewInt(0))
	if err != nil {
		t.Fatal("Get root public key fail: ", err)
	}

	if !pub.Verify(rb, sig) {
		t.Fatal("Verify root signature fail: ", err)
	}

	// Child Private Key Signing

	index := big.NewInt(0x1000)

	sig2, err := crypto.PrivateKeySign(priv, index, rb)
	if err != nil {
		t.Fatal("Sign raw data fail: ", err)
	}

	// root private -> child private -> child public
	pub2, err := priv.ChildPublic(index)
	if err != nil {
		t.Fatal("Get child public key fail: ", err)
	}

	if !pub2.Verify(rb, sig2) {
		t.Fatal("Verify child signature fail: ", err)
	}

	// root private -> root public -> child public
	pub3, err := pub.child(index)
	if err != nil {
		t.Fatal("Get child public key fail: ", err)
	}

	if !pub3.Verify(rb, sig2) {
		t.Fatal("Verify child signature fail: ", err)
	}

	//t.Logf("Private Key: %v", priv)
	//t.Logf("Public Key: %v", pub)
}

func TestPrivateKey_Serialize(t *testing.T) {

	// Generate private key
	privateKeySerialize(t)

}

func BenchmarkPrivateKey_ECP256_FIPS186_Generate(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_, err := NewPrivatekey(ECP256_FIPS186)
		if err != nil {
			b.Fatal("Generate private key fail: ", err)
		}
	}
}

func BenchmarkPrivateKey_SECP256K1_Generate(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_, err := NewPrivatekey(SECP256K1)
		if err != nil {
			b.Fatal("Generate private key fail: ", err)
		}
	}
}

func BenchmarkPrivateKey_ECP256_FIPS186_Sign(b *testing.B) {

	priv, err := NewPrivatekey(ECP256_FIPS186)
	if err != nil {
		b.Fatal("Generate private key fail: ", err)
	}

	// Generate message data
	rb := make([]byte, 32)
	_, err = rand.Read(rb)
	if err != nil {
		b.Fatal("Genereate rand bytes fail: ", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Signing
		_, err := priv.Sign(rb)
		if err != nil {
			b.Fatal("Sign raw data fail: ", err)
		}
	}
}

func BenchmarkPrivateKey_SECP256K1_Sign(b *testing.B) {

	priv, err := NewPrivatekey(SECP256K1)
	if err != nil {
		b.Fatal("Generate private key fail: ", err)
	}

	// Generate message data
	rb := make([]byte, 32)
	_, err = rand.Read(rb)
	if err != nil {
		b.Fatal("Genereate rand bytes fail: ", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Signing
		_, err := priv.Sign(rb)
		if err != nil {
			b.Fatal("Sign raw data fail: ", err)
		}
	}
}

func BenchmarkPrivateKey_ECP256_FIPS186_Verify(b *testing.B) {

	priv, err := NewPrivatekey(ECP256_FIPS186)
	if err != nil {
		b.Fatal("Generate private key fail: ", err)
	}

	// Generate message data
	rb := make([]byte, 32)
	_, err = rand.Read(rb)
	if err != nil {
		b.Fatal("Genereate rand bytes fail: ", err)
	}

	// Signing
	sig, err := priv.Sign(rb)
	if err != nil {
		b.Fatal("Sign raw data fail: ", err)
	}

	// Public key
	pub := priv.Public()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if !pub.Verify(rb, sig) {
			b.Fatal("Verify root signature fail: ", err)
		}
	}
}

func BenchmarkPrivateKey_SECP256K1_Verify(b *testing.B) {

	priv, err := NewPrivatekey(SECP256K1)
	if err != nil {
		b.Fatal("Generate private key fail: ", err)
	}

	// Generate message data
	rb := make([]byte, 32)
	_, err = rand.Read(rb)
	if err != nil {
		b.Fatal("Genereate rand bytes fail: ", err)
	}

	// Signing
	sig, err := priv.Sign(rb)
	if err != nil {
		b.Fatal("Sign raw data fail: ", err)
	}

	// Public key
	pub := priv.Public()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if !pub.Verify(rb, sig) {
			b.Fatal("Verify root signature fail: ", err)
		}
	}
}

func BenchmarkPrivateKey_Serialize(b *testing.B) {

	for i := 0; i < b.N; i++ {
		privateKeySerialize(b)
	}
}
func privateKeySerialize(t testing.TB) {

	// Generate private key

	priv, err := NewPrivatekey(DefaultCurveType)
	if err != nil {
		t.Fatal("Generate private key fail: ", err)
	}

	// Dump private key

	privPB := priv.PBMessage()

	// Import private key

	priv2 := new(PrivateKey)
	if err := priv2.FromPBMessage(privPB); err != nil {
		t.Fatal("Import private key fail: ", err)
	}

	// Compare

	if !priv.IsEqualForTest(priv2) {
		t.Fatal("Private keys not equal")
	}
}
