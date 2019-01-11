package ecdsa

import (
	"crypto/ecdsa"
	"crypto/rand"
	"hyperledger.abchain.org/core/crypto"
	"math/big"
	"testing"
)

func genPrivkey(eci int32, s string) *ecdsa.PrivateKey {

	ec, err := GetEC(eci)
	if err != nil {
		panic("ec fail")
	}

	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = ec
	priv.D, _ = big.NewInt(0).SetString(s, 10)
	priv.PublicKey.X, priv.PublicKey.Y = ec.ScalarBaseMult(priv.D.Bytes())

	return priv
}

func TestSECP256K1Correction(t *testing.T) {

	priv := genPrivkey(SECP256K1, "110871560887871433772251872457573575186844291339401824426745393239615136047255")
	sigs, _ := big.NewInt(0).SetString("50551231335931353095826110317475557373687171680522884359116420974603396918366", 10)
	sigr, _ := big.NewInt(0).SetString("6861113568642063833136582234362472053402008863274491262222906151310795201436", 10)

	pubX := big.NewInt(0).SetBytes([]byte{0x28, 0x5c, 0x55, 0x01, 0xa7, 0x08, 0x5b, 0xa0, 0xbb, 0x4b, 0xf7, 0x1e, 0x2d, 0x46, 0x4f, 0xf1, 0xae, 0x69, 0x1a, 0x3f, 0xaa, 0x89, 0xee, 0x31, 0x2a, 0xbc, 0xc9, 0x1d, 0xad, 0x21, 0xb2, 0x81})
	if priv.PublicKey.X.Cmp(pubX) != 0 {
		t.Fatalf("unexpected pubX (%s, expect %s)", priv.PublicKey.X, pubX)
	}

	sigHash := []byte{0x0c, 0x28, 0xfc, 0xa3, 0x86, 0xc7, 0xa2, 0x27, 0x60, 0x0b, 0x2f, 0xe5, 0x0b, 0x7c, 0xae, 0x11, 0xec, 0x86, 0xd3, 0xbf, 0x1f, 0xbe, 0x47, 0x1b, 0xe8, 0x98, 0x27, 0xe1, 0x9d, 0x72, 0xaa, 0x1d}
	if !ecdsa.Verify(&priv.PublicKey, sigHash, sigr, sigs) {
		t.Fatal("Verify sign failure")
	}

	priv2 := genPrivkey(SECP256K1, "70912661286874047882889983861262634846432069668246130355602184839087556488003")
	sigs, _ = big.NewInt(0).SetString("53428889315552072053467946166248932714858581566670823688759633209466837559593", 10)
	sigr, _ = big.NewInt(0).SetString("8698353116727175722979305956484862846299506136674989554075290215293275990614", 10)

	pubX = big.NewInt(0).SetBytes([]byte{0x89, 0x48, 0x85, 0x14, 0xb6, 0x4b, 0xed, 0x5a, 0x13, 0x9d, 0xf3, 0xfb, 0x97, 0xc3, 0xde, 0xdb, 0xe5, 0x12, 0x64, 0x4a, 0x6e, 0x66, 0xc3, 0x07, 0xa0, 0x99, 0xaf, 0xcd, 0xd2, 0xf7, 0x5b, 0x3a})
	if priv2.PublicKey.X.Cmp(pubX) != 0 {
		t.Fatalf("unexpected pubX (%s, expect %s)", priv2.PublicKey.X, pubX)
	}

	if !ecdsa.Verify(&priv2.PublicKey, sigHash, sigr, sigs) {
		t.Fatal("Verify sign 2 failure")
	}

}

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

	if sig.Kd != nil {
		t.Fatal("Sign data include ghost derived data", sig)
	}

	pub := priv.public()

	if !pub.Verify(rb, sig) {
		t.Fatal("Verify root signature fail: ", err)
	}

	// Child Private Key Signing

	index := big.NewInt(0x1000)

	sig2, err := crypto.PrivateKeySign(priv, index, rb)
	if err != nil {
		t.Fatal("Sign raw data fail: ", err)
	}

	if sig2.Kd == nil {
		t.Fatal("Sign data not include derived data")
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

func TestPubRecover(t *testing.T) {

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

	pub3 := new(PublicKey)
	err = pub3.Recover(sig2)
	if err != nil {
		t.Fatal("Recover public key fail: ", err)
	}

	if !pub2.Verify(rb, sig2) {
		t.Fatal("Verify child signature fail: ", err)
	}

	if !pub2.IsEqualForTest(pub3) {
		t.Fatal("recover pubkey is not total equal", pub2, pub3)
	}

	//deprivate derived data
	sig2.Kd = nil
	err = pub3.Recover(sig2)
	if err != nil {
		t.Fatal("Recover public key fail: ", err)
	}

	if !pub2.IsEqual(pub3) {
		t.Fatal("recover pubkey is not equal", pub2, pub3)
	}

	if pub2.IsEqualForTest(pub3) {
		t.Fatal("recover pubkey has ghost derived data", pub2, pub3)
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
