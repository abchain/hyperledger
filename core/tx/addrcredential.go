package abchainTx

import (
	"bytes"
	"errors"
	_ "fmt"
	"hyperledger.abchain.org/core/crypto"
	_ "hyperledger.abchain.org/core/utils"
	pb "hyperledger.abchain.org/protos"
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

type addrCred struct {
	c            *pb.TxCredential_AddrCredentials
	verified     bool
	verifiedAddr Address
}

func (addrc *addrCred) GetCredPubkey(addr Address) *crypto.PublicKey {

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

func (addrc *addrCred) Verify(addr Address, hash []byte) error {

	switch v := addrc.c.Cred.(type) {
	case *pb.TxCredential_AddrCredentials_User:
		if addrc.verified {
			if addrc.verifiedAddr.IsEqual(&addr) {
				return nil
			} else {
				return errors.New("address is not matched with verified one")
			}
		}

		pk := addrc.GetCredPubkey(addr)
		if pk == nil {
			return errors.New("address is not matched with cred")
		}

		if err := verifyPk(pk, hash, v.User); err != nil {
			return err
		}

		addrc.verified = true
		addrc.verifiedAddr = addr
		return nil
	case *pb.TxCredential_AddrCredentials_Cc:
		return errors.New("cc calling has been deprecated")
	default:
		return errors.New("Cred type not recongnized")
	}
}

type soleAddrCred struct {
	hash []byte
	*addrCred
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
		return nil, errors.New("use deprecatd cc credential")
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

const (
	fingerprintIndexLength int = 4
)

func (addrc *soleAddrCred) Verify(addr Address) error {
	return addrc.addrCred.Verify(addr, addrc.hash)
}

func (addrc *soleAddrCred) CredCount() int {
	return 1
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
	creds map[uint32]*addrCred
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

	return cred.Verify(addr, maddrc.hash)
}

func (maddrc *mutipleAddrCred) CredCount() int {
	return len(maddrc.creds)
}

func (maddrc *mutipleAddrCred) GetCredPubkey(addr Address) *crypto.PublicKey {

	c, ok := maddrc.creds[toFingerPrint(addr.Hash)]

	if !ok {
		return nil
	}

	pk, err := pubkeyFromCred(c.c)
	if err != nil {
		return nil
	}

	return pk
}

func (maddrc *mutipleAddrCred) ListCredPubkeys() []*crypto.PublicKey {

	ret := make([]*crypto.PublicKey, 0, maddrc.CredCount())

	for _, c := range maddrc.creds {
		pk, _ := pubkeyFromCred(c.c)
		if pk != nil {
			ret = append(ret, pk)
		}
	}

	return ret
}

func NewAddrCredential(hash []byte, addrc []*pb.TxCredential_AddrCredentials) (ret AddrCredentials, e error) {

	switch len(addrc) {
	case 0:
		ret = &noAddrCred{}
	case 1:
		ret = &soleAddrCred{hash, &addrCred{c: addrc[0]}}
	default:

		rret := &mutipleAddrCred{hash, make(map[uint32]*addrCred)}
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

			rret.creds[toFingerPrint(fp)] = &addrCred{c: c}
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

	b.e = errors.New("cc credential has been deprecated")
	// if len(addr.Hash) < b.indexLen {

	// 	return
	// }

	// b.cache[toFingerPrint(addr.Hash)] = &builderCache{
	// 	addr.Hash,
	// 	pb.TxCredential_AddrCredentials{
	// 		&pb.TxCredential_AddrCredentials_Cc{
	// 			&pb.TxCredential_InnerCredential{
	// 				ccname,
	// 				nil,
	// 				pub.PBMessage(),
	// 			},
	// 		},
	// 	},
	// }
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
