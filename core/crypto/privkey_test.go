package crypto

import (
	"crypto/rand"
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

	sig, err := priv.Sign(big.NewInt(0), rand.Reader, rb)
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

	sig2, err := priv.Sign(index, rand.Reader, rb)
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
	pub3, err := pub.ChildKey(index)
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

	priv, err := NewPrivatekey(DefaultCurveType)
	if err != nil {
		t.Fatal("Generate private key fail: ", err)
	}

	// Dump private key

	privStr := priv.Str()
	t.Logf("Dump private key: %s", privStr)

	// Import private key

	priv2, err := PrivatekeyFromString(privStr)
	if err != nil {
		t.Fatal("Import private key fail: ", err)
	}

	// Compare

	if !priv.IsEqualForTest(priv2) {
		t.Fatal("Private keys not equal")
	}
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
		_, err := priv.Sign(big.NewInt(0), rand.Reader, rb)
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
		_, err := priv.Sign(big.NewInt(0), rand.Reader, rb)
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
	sig, err := priv.Sign(big.NewInt(0), rand.Reader, rb)
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
	sig, err := priv.Sign(big.NewInt(0), rand.Reader, rb)
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
		privateKeySerializeBenchmark(b)
	}
}

func privateKeySerializeBenchmark(b *testing.B) {

	// Generate private key
	priv, err := NewPrivatekey(DefaultCurveType)
	if err != nil {
		b.Fatal("Generate private key fail: ", err)
	}

	// Dump private key
	privStr := priv.Str()

	// Import private key
	priv2, err := PrivatekeyFromString(privStr)
	if err != nil {
		b.Fatal("Import private key fail: ", err)
	}

	// Compare
	if !priv.IsEqualForTest(priv2) {
		b.Fatal("Private keys not equal")
	}
}
