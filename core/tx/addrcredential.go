package abchainTx

import (
	"errors"

	"hyperledger.abchain.org/core/crypto"
	pb "hyperledger.abchain.org/protos"
)

type noAddrCred struct {
}

func (c *noAddrCred) Verify(Address) error {
	return errors.New("No crendential")
}

func (c *noAddrCred) CredCount() int {
	return 0
}

func (c *noAddrCred) GetCredPubkey(addr Address) crypto.Verifier {
	return nil
}

func (c *noAddrCred) ListCredPubkeys() []crypto.Verifier {
	return nil
}

type addrCred struct {
	c            *pb.TxCredential_AddrCredentials
	verified     bool
	verifiedAddr Address
}

func (addrc *addrCred) GetCredPubkey(addr Address) crypto.Verifier {

	pk, err := pubkeyFromCred(addrc.c)
	if err != nil {
		return nil
	}

	innerAddr, err := NewAddress(pk)
	if err != nil {
		return nil
	}

	if !addr.IsEqual(innerAddr) {
		return nil
	}

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
		pk, err := crypto.PublicKeyFromSignature(v.User.Signature)

		if err != nil {
			return nil, err
		}

		return pk.Digest(), nil
	case *pb.TxCredential_AddrCredentials_Cc:
		return v.Cc.Fingerprint, nil
	default:
		return nil, errors.New("Unrecognized cred type")
	}
}

func pubkeyFromCred(c *pb.TxCredential_AddrCredentials) (crypto.Verifier, error) {

	switch v := c.Cred.(type) {
	case *pb.TxCredential_AddrCredentials_User:
		return crypto.PublicKeyFromSignature(v.User.GetSignature())
	case *pb.TxCredential_AddrCredentials_Cc:
		return nil, errors.New("use deprecatd cc credential")
	default:
		return nil, errors.New("Unrecognized cred type")
	}

}

func verifyPk(pk crypto.Verifier, hash []byte, u *pb.TxCredential_UserCredential) error {

	if u.Signature == nil {
		return errors.New("No signature")
	}

	if pk.Verify(hash, u.Signature) {
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

func (addrc *soleAddrCred) ListCredPubkeys() []crypto.Verifier {

	pk, err := pubkeyFromCred(addrc.c)
	if err != nil {
		return nil
	}

	return []crypto.Verifier{pk}
}

type mutipleAddrCred struct {
	hash  []byte
	creds map[uint32]*addrCred
}

var powconst = [4]uint32{1, 256, 65536, 1677216}

func toFingerPrint(b []byte) (ret uint32) {

	var blen = fingerprintIndexLength
	if len(b) < blen {
		blen = len(b)
	}

	for i := 0; i < blen; i++ {
		ret += uint32(b[i]) * powconst[i]
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

func (maddrc *mutipleAddrCred) GetCredPubkey(addr Address) crypto.Verifier {

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

func (maddrc *mutipleAddrCred) ListCredPubkeys() []crypto.Verifier {

	ret := make([]crypto.Verifier, 0, maddrc.CredCount())

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

func (b *builder) AddSignature(sign *pb.Signature) {

	pk, err := crypto.PublicKeyFromSignature(sign)
	if err != nil {
		b.e = err
		return
	}

	bt := pk.Digest()
	if len(bt) < b.indexLen {
		b.e = errors.New("Wrong index for pk: too short")
		return
	}

	ind := toFingerPrint(bt)
	if _, ok := b.cache[ind]; ok {
		b.e = errors.New("Duplicated signature from same source")
		return
	}

	b.cache[ind] = &builderCache{bt,
		pb.TxCredential_AddrCredentials{
			&pb.TxCredential_AddrCredentials_User{&pb.TxCredential_UserCredential{sign}},
		},
	}
}

func (b *builder) AddCc(ccname string, addr Address, pub crypto.Verifier) {

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
