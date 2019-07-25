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
	*pb.TxCredential_AddrCredentials_User
	verified     bool
	verifiedAddr Address
}

func (addrc *addrCred) GetCredPubkey(addr Address) crypto.Verifier {

	pk, err := crypto.PublicKeyFromSignature(addrc.User.GetSignature())
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

	if err := verifyPk(pk, hash, addrc.User); err != nil {
		return err
	}

	addrc.verified = true
	addrc.verifiedAddr = addr
	return nil

}

type soleAddrCred struct {
	hash []byte
	*addrCred
}

func addrFingerPrintFromCred(v *pb.TxCredential_AddrCredentials_User) ([]byte, error) {

	pk, err := crypto.PublicKeyFromSignature(v.User.Signature)

	if err != nil {
		return nil, err
	}

	return pk.Digest(), nil
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

	pk, err := crypto.PublicKeyFromSignature(addrc.User.GetSignature())
	if err != nil {
		return nil
	}

	return []crypto.Verifier{pk}
}

type mutipleAddrCred struct {
	hash  []byte
	creds map[uint32][]*addrCred
	size  int
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

func (maddrc *mutipleAddrCred) Verify(addr Address) (err error) {

	creds, ok := maddrc.creds[toFingerPrint(addr.Hash)]

	if !ok {
		return errors.New("no cred for pk addr")
	}

	for _, cred := range creds {
		if err = cred.Verify(addr, maddrc.hash); err == nil {
			return nil
		}
	}

	return
}

func (maddrc *mutipleAddrCred) CredCount() int {
	return maddrc.size
}

func (maddrc *mutipleAddrCred) GetCredPubkey(addr Address) crypto.Verifier {

	cs, ok := maddrc.creds[toFingerPrint(addr.Hash)]

	if !ok {
		return nil
	}

	for _, c := range cs {
		pk := c.GetCredPubkey(addr)
		if pk != nil {
			return pk
		}
	}
	return nil
}

func (maddrc *mutipleAddrCred) ListCredPubkeys() []crypto.Verifier {

	ret := make([]crypto.Verifier, 0, maddrc.CredCount())

	for _, cs := range maddrc.creds {
		for _, c := range cs {
			pk, _ := crypto.PublicKeyFromSignature(c.User.GetSignature())
			if pk != nil {
				ret = append(ret, pk)
			}
		}
	}

	return ret
}

func NewAddrCredential(hash []byte, addrc []*pb.TxCredential_AddrCredentials_User) (ret AddrCredentials, e error) {

	switch len(addrc) {
	case 0:
		ret = &noAddrCred{}
	case 1:
		ret = &soleAddrCred{hash, &addrCred{TxCredential_AddrCredentials_User: addrc[0]}}
	default:

		rret := &mutipleAddrCred{hash, make(map[uint32][]*addrCred), len(addrc)}
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

			ind := toFingerPrint(fp)
			rret.creds[ind] = append(rret.creds[ind],
				&addrCred{TxCredential_AddrCredentials_User: c})
		}
	}

	return
}
