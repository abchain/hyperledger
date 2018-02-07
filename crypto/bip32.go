package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/sha512"
	"math/big"
)

func getChildPrivateKey(root *PrivateKey, index *big.Int) (*PrivateKey, error) {

	curve, err := GetEC(root.CurveType)
	if err != nil {
		return nil, err
	}

	// Use pubkey to generate intermediary
	data := root.Key.PublicKey.X.Bytes()
	data = append(data, root.Key.PublicKey.Y.Bytes()...)
	intermediary, err := getIntermediary(data, root.Chaincode, index)
	if err != nil {
		return nil, err
	}

	childD := addPrivateKeys(root.Key.D.Bytes(), intermediary[:32], curve)
	if err = validatePrivateKey(childD, curve.Params()); err != nil {
		return nil, err
	}

	x, y := curve.ScalarBaseMult(childD)
	childPubkey := &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}

	if err = validateChildPublicKey(childPubkey); err != nil {
		return nil, err
	}

	childECDSAKey := &ecdsa.PrivateKey{
		PublicKey: *childPubkey,
		D:         new(big.Int).SetBytes(childD),
	}

	rootFingerprint, err := GetPrivateKeyRootFingerprint(root.Key)
	if err != nil {
		return nil, err
	}

	child := &PrivateKey{
		root.Version,
		root.CurveType,
		childECDSAKey,
		rootFingerprint,
		index,
		intermediary[32:],
	}

	return child, nil
}

func getChildPublicKey(root *PublicKey, index *big.Int) (*PublicKey, error) {

	curve, err := GetEC(root.CurveType)
	if err != nil {
		return nil, err
	}

	// Use pubkey to generate intermediary
	data := root.Key.X.Bytes()
	data = append(data, root.Key.Y.Bytes()...)
	intermediary, err := getIntermediary(data, root.Chaincode, index)
	if err != nil {
		return nil, err
	}

	x, y := curve.ScalarBaseMult(intermediary[:32])
	intermediaryPubkey := &ecdsa.PublicKey{curve, x, y}

	childKey := addPublicKeys(intermediaryPubkey, root.Key, curve)

	if err = validateChildPublicKey(childKey); err != nil {
		return nil, err
	}

	rootFingerprint, err := GetPublicKeyRootFingerprint(root.Key)
	if err != nil {
		return nil, err
	}

	child := &PublicKey{
		root.Version,
		root.CurveType,
		childKey,
		rootFingerprint,
		index,
		intermediary[32:],
	}

	return child, nil
}

func getIntermediary(key []byte, chaincode []byte, index *big.Int) ([]byte, error) {
	data := key
	data = append(data, index.Bytes()...)

	hmac := hmac.New(sha512.New, chaincode)
	_, err := hmac.Write(data)
	if err != nil {
		return nil, err
	}

	return hmac.Sum(nil), nil
}

func addPublicKeys(key1 *ecdsa.PublicKey, key2 *ecdsa.PublicKey,
	curve elliptic.Curve) *ecdsa.PublicKey {

	X, Y := curve.Add(key1.X, key1.Y, key2.X, key2.Y)

	return &ecdsa.PublicKey{
		curve,
		X,
		Y,
	}
}

func addPrivateKeys(key1 []byte, key2 []byte, curve elliptic.Curve) []byte {
	var key1Int big.Int
	var key2Int big.Int
	key1Int.SetBytes(key1)
	key2Int.SetBytes(key2)

	key1Int.Add(&key1Int, &key2Int)
	key1Int.Mod(&key1Int, curve.Params().N)

	b := key1Int.Bytes()
	if len(b) < 32 {
		extra := make([]byte, 32-len(b))
		b = append(extra, b...)
	}
	return b
}

func validateChildPublicKey(key *ecdsa.PublicKey) error {
	if key.X.Sign() == 0 || key.Y.Sign() == 0 {
		return ErrInvalidPublicKey
	}

	return nil
}
