package abchainTx

import (
	"bytes"
	"errors"
	_ "fmt"
	"github.com/abchain/fabric/core/chaincode/shim"
	"hyperledger.abchain.org/crypto"
	pb "hyperledger.abchain.org/protos"
	_ "hyperledger.abchain.org/utils"
	_ "math/big"
)

type noAddrCred struct {
}

func (c *noAddrCred) Verify(Address) error {
	return errors.New("No crendential")
}

func (c *noAddrCred) CredCount() int {
	return 0
}

func (c *noAddrCred) GetCredPubkey(addr Address) *crypto.PublicKey {
	return nil
}

func (c *noAddrCred) ListCredPubkeys() []*crypto.PublicKey {
	return nil
}

type soleAddrCred struct {
	hash []byte
	stub shim.ChaincodeStubInterface
	c    *pb.TxCredential_AddrCredentials
}

func addrFingerPrintFromCred(c *pb.TxCredential_AddrCredentials) ([]byte, error) {
	switch v := c.Cred.(type) {
	case *pb.TxCredential_AddrCredentials_User:
		pk, err := crypto.PublicKeyFromPBMessage(v.User.Pk)

		if err != nil {
			return nil, err
		}

		pkh, err := GetPublicKeyHash(pk)

		if err != nil {
			return nil, err
		}

		return pkh, nil
	case *pb.TxCredential_AddrCredentials_Cc:
		return v.Cc.Fingerprint, nil
	default:
		return nil, errors.New("Unrecognized cred type")
	}
}

func pubkeyFromCred(c *pb.TxCredential_AddrCredentials) (*crypto.PublicKey, error) {

	var pkk *pb.PublicKey

	switch v := c.Cred.(type) {
	case *pb.TxCredential_AddrCredentials_User:
		pkk = v.User.Pk
	case *pb.TxCredential_AddrCredentials_Cc:
		pkk = v.Cc.Pk
	default:
		return nil, errors.New("Unrecognized cred type")
	}

	return crypto.PublicKeyFromPBMessage(pkk)
}

func verifyPk(pk *crypto.PublicKey, hash []byte, u *pb.TxCredential_UserCredential) error {

	sign, err := crypto.ECSignatureFromPB(u.Signature)

	if err != nil {
		return err
	}

	if sign.Verify(hash, pk) {
		return nil
	}

	return errors.New("signature is not matched with pubkey")

}

func verifyCc(u *pb.TxCredential_InnerCredential, _ shim.ChaincodeStubInterface) error {

	//TODO: should verify ccname with stub
	return errors.New("cc calling has no implement yet")
}

const (
	fingerprintIndexLength int = 4
)

func (addrc *soleAddrCred) Verify(addr Address) error {

	pk := addrc.GetCredPubkey(addr)
	if pk == nil {
		return errors.New("address is not matched with cred")
	}

	switch v := addrc.c.Cred.(type) {
	case *pb.TxCredential_AddrCredentials_User:
		return verifyPk(pk, addrc.hash, v.User)
	case *pb.TxCredential_AddrCredentials_Cc:
		return verifyCc(v.Cc, addrc.stub)
	}

	return errors.New("Cred type not recongnized")
}

func (addrc *soleAddrCred) CredCount() int {
	return 1
}

func (addrc *soleAddrCred) GetCredPubkey(addr Address) *crypto.PublicKey {

	fp, err := addrFingerPrintFromCred(addrc.c)
	if err != nil {
		return nil
	}

	if len(fp) < fingerprintIndexLength {
		return nil
	}

	if bytes.Compare(addr.Hash[:len(fp)], fp) != 0 {
		return nil
	}

	pk, err := pubkeyFromCred(addrc.c)
	return pk

}

func (addrc *soleAddrCred) ListCredPubkeys() []*crypto.PublicKey {

	pk, err := pubkeyFromCred(addrc.c)
	if err != nil {
		return nil
	}

	return []*crypto.PublicKey{pk}
}

type mutipleAddrCred struct {
	hash  []byte
	stub  shim.ChaincodeStubInterface
	creds map[uint32]*pb.TxCredential_AddrCredentials
}

func toFingerPrint(b []byte) (ret uint32) {

	var blen = fingerprintIndexLength
	if len(b) < blen {
		blen = len(b)
	}

	pow := [4]uint32{1, 256, 65536, 1677216}

	for i := 0; i < blen; i++ {
		ret += uint32(b[i]) * pow[i]
	}

	return
}

func (maddrc *mutipleAddrCred) Verify(addr Address) error {

	cred, ok := maddrc.creds[toFingerPrint(addr.Hash)]

	if !ok {
		return errors.New("no cred for pk addr")
	}

	tc := &soleAddrCred{
		maddrc.hash,
		maddrc.stub,
		cred,
	}

	return tc.Verify(addr)
}

func (maddrc *mutipleAddrCred) CredCount() int {
	return len(maddrc.creds)
}

func (maddrc *mutipleAddrCred) GetCredPubkey(addr Address) *crypto.PublicKey {

	c, ok := maddrc.creds[toFingerPrint(addr.Hash)]

	if !ok {
		return nil
	}

	pk, err := pubkeyFromCred(c)
	if err != nil {
		return nil
	}

	return pk
}

func (maddrc *mutipleAddrCred) ListCredPubkeys() []*crypto.PublicKey {

	ret := make([]*crypto.PublicKey, 0, maddrc.CredCount())

	for _, c := range maddrc.creds {
		pk, _ := pubkeyFromCred(c)
		if pk != nil {
			ret = append(ret, pk)
		}
	}

	return ret
}

func NewAddrCredential(hash []byte, stub shim.ChaincodeStubInterface,
	addrc []*pb.TxCredential_AddrCredentials) (ret AddrCredentials, e error) {

	switch len(addrc) {
	case 0:
		ret = &noAddrCred{}
	case 1:
		ret = &soleAddrCred{hash, stub, addrc[0]}
	default:

		rret := &mutipleAddrCred{hash, stub, make(map[uint32]*pb.TxCredential_AddrCredentials)}
		ret = rret
		for _, c := range addrc {

			fp, err := addrFingerPrintFromCred(c)
			if err != nil {
				e = err
				return
			}

			if len(fp) < fingerprintIndexLength {
				e = errors.New("Invalid fingerprint: too short")
				return
			}

			rret.creds[toFingerPrint(fp)] = c
		}
	}

	return
}

type builderCache struct {
	index []byte
	d     pb.TxCredential_AddrCredentials
}

type builder struct {
	e        error
	indexLen int
	cache    map[uint32]*builderCache
}

func NewAddrCredentialBuilder() AddrCredentialBuilder {
	return &builder{nil, fingerprintIndexLength, make(map[uint32]*builderCache)}
}

func (b *builder) AddSignature(pk *crypto.PublicKey, sign *crypto.ECSignature) {

	bt, err := GetPublicKeyHash(pk)
	if err != nil {
		b.e = err
		return
	}

	if len(bt) < b.indexLen {
		b.e = errors.New("Wrong index for pk: too short")
		return
	}

	cache := &builderCache{bt,
		pb.TxCredential_AddrCredentials{
			&pb.TxCredential_AddrCredentials_User{
				&pb.TxCredential_UserCredential{
					sign.ToPB(),
					pk.PBMessage(),
				},
			},
		},
	}

	b.cache[toFingerPrint(bt)] = cache
}

func (b *builder) AddCc(ccname string, addr Address, pub *crypto.PublicKey) {

	if len(addr.Hash) < b.indexLen {
		b.e = errors.New("Wrong index for pk: too short")
		return
	}

	b.cache[toFingerPrint(addr.Hash)] = &builderCache{
		addr.Hash,
		pb.TxCredential_AddrCredentials{
			&pb.TxCredential_AddrCredentials_Cc{
				&pb.TxCredential_InnerCredential{
					ccname,
					nil,
					pub.PBMessage(),
				},
			},
		},
	}
}

func (b *builder) Update(msg *pb.TxCredential) error {

	if b.e != nil {
		return b.e
	}

	msg.Addrc = make([]*pb.TxCredential_AddrCredentials, 0, len(b.cache))

	for _, c := range b.cache {
		switch v := c.d.Cred.(type) {
		case *pb.TxCredential_AddrCredentials_Cc:
			v.Cc.Fingerprint = c.index[:b.indexLen]
		}

		msg.Addrc = append(msg.Addrc, &c.d)
	}

	return nil
}
