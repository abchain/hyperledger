package ecdsa

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/sha512"
	"errors"
	"hyperledger.abchain.org/protos"
	"math/big"
)

type KeyDerivation struct {
	// underlying root private key fingerprint
	RootFingerPrint []byte
	// index of child private key
	Index *big.Int
	// chaincode
	Chaincode []byte
}

func (kd *KeyDerivation) GetIndex() *big.Int {
	if kd == nil {
		return big.NewInt(0)
	}
	return kd.Index
}

func (kd *KeyDerivation) GetRootFingerPrint() []byte {

	if kd == nil {
		return nil
	}

	return kd.RootFingerPrint
}

func (kd *KeyDerivation) ToKDMessage() *protos.KeyDerived {

	if kd == nil {
		return nil
	}

	return &protos.KeyDerived{
		RootFingerprint: kd.RootFingerPrint,
		Index:           kd.Index.Bytes(),
		Chaincode:       kd.Chaincode,
	}
}

func (kd *KeyDerivation) isEqual(otherkd *KeyDerivation) bool {

	if kd == nil {
		return otherkd == nil
	} else if otherkd == nil {
		return false
	}

	if !bytes.Equal(kd.Chaincode, otherkd.Chaincode) {
		return false
	}

	if !bytes.Equal(kd.RootFingerPrint, otherkd.RootFingerPrint) {
		return false
	}

	if kd.Index.Cmp(otherkd.Index) != 0 {
		return false
	}

	return true
}

func (kd *KeyDerivation) FromKDMessage(msg *protos.KeyDerived) {
	kd.RootFingerPrint = msg.GetRootFingerprint()
	kd.Index = big.NewInt(0).SetBytes(msg.GetIndex())
	kd.Chaincode = msg.GetChaincode()
}

func paddingBytes(bt []byte, size int) []byte {
	if len(bt) > size {
		panic("Large bytes")
	} else if len(bt) == size {
		return bt
	} else {
		return append(make([]byte, size-len(bt)), bt...)
	}
}

func (kd *KeyDerivation) GenIntermediary(pub *ecdsa.PublicKey, index *big.Int) ([]byte, *KeyDerivation, error) {

	if kd == nil {
		return nil, nil, errors.New("No keyderivation data")
	} else if len(kd.Chaincode) < 32 {
		//we require 256 bit cc at least
		return nil, nil, errors.New("Not valid chaincode")
	}

	cmpPoint := compressPublicKey(pub.X, pub.Y)
	indBytes := index.Bytes()
	if len(indBytes) < 4 {
		indBytes = append(make([]byte, 4-len(indBytes)), indBytes...)
	}

	data := bytes.Join([][]byte{cmpPoint, indBytes}, nil)

	hmac := hmac.New(sha512.New, kd.Chaincode)
	_, err := hmac.Write(data)
	if err != nil {
		return nil, nil, err
	}

	rootFingerprint, err := digestECDSAPk(pub)
	if err != nil {
		return nil, nil, err
	}

	//sanity check
	if len(rootFingerprint) < PUBLICKEY_FINGERPRINT_LEN {
		panic("Wrong private key fingerprint length")
	}

	ddata := hmac.Sum(nil)
	if len(ddata) < 64 {
		panic("HMAC in 512bit hash give no enough bits")
	}

	return ddata[:32], &KeyDerivation{rootFingerprint[:PUBLICKEY_FINGERPRINT_LEN], index, ddata[32:]}, nil
}

func getChildPrivateKey(root *PrivateKey, index *big.Int) (*PrivateKey, error) {

	curve, err := GetEC(root.CurveType)
	if err != nil {
		return nil, err
	}

	// Use pubkey to generate intermediary
	drvkey, drvkd, err := root.GenIntermediary(&root.Key.PublicKey, index)
	if err != nil {
		return nil, err
	}

	childD := addPrivateKeys(root.Key.D.Bytes(), drvkey, curve)
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

	return &PrivateKey{
		Version:       root.Version,
		CurveType:     root.CurveType,
		KeyDerivation: drvkd,
		Key: &ecdsa.PrivateKey{
			PublicKey: *childPubkey,
			D:         big.NewInt(0).SetBytes(childD),
		},
	}, nil
}

func getChildPublicKey(root *PublicKey, index *big.Int) (*PublicKey, error) {

	curve, err := GetEC(root.CurveType)
	if err != nil {
		return nil, err
	}

	drvkey, drvkd, err := root.GenIntermediary(root.Key, index)
	if err != nil {
		return nil, err
	}

	x, y := curve.ScalarBaseMult(drvkey)

	childKey := addPublicKeys(&ecdsa.PublicKey{curve, x, y}, root.Key, curve)

	if err = validateChildPublicKey(childKey); err != nil {
		return nil, err
	}

	return &PublicKey{
		Version:       root.Version,
		CurveType:     root.CurveType,
		KeyDerivation: drvkd,
		Key:           childKey,
	}, nil
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
		return errors.New("Invalid derived ECDSA pubkey")
	}

	return nil
}

func validatePrivateKey(d []byte, curveParams *elliptic.CurveParams) error {

	dint := big.NewInt(0).SetBytes(d)

	if dint.Int64() == 0 || dint.Cmp(curveParams.N) >= 0 {
		return errors.New("Invalid derived ECDSA privatekey")
	}
	return nil
}
