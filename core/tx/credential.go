package abchainTx

import (
	"fmt"
	"hyperledger.abchain.org/core/crypto"
	pb "hyperledger.abchain.org/protos"
)

type AddrCredentials interface {
	Verify(addr Address) error
	CredCount() int
	// GetCredPubkey(addr Address) *crypto.PublicKey
	// ListCredPubkeys() []*crypto.PublicKey
	GetCredPubkey(addr Address) crypto.Verifier
	ListCredPubkeys() []crypto.Verifier
}

type DataCredentials interface {
	GetValue(string) interface{}
}

type txCredentials struct {
	AddrCredentials
	DataCredentials
}

func newTxCredential(hash []byte, cred []*pb.TxCredential_AddrCredentials) (txCredentials, error) {

	var ac []*pb.TxCredential_AddrCredentials_User
	dataCred := dataCred(make(map[string]interface{}))

	for _, c := range cred {
		switch v := c.Cred.(type) {
		case *pb.TxCredential_AddrCredentials_User:
			ac = append(ac, v)
		case *pb.TxCredential_AddrCredentials_Data:
			switch vv := v.Data.GetData().(type) {
			case *pb.TxCredential_DataCredential_Bts:
				dataCred[v.Data.GetKey()] = vv.Bts
			case *pb.TxCredential_DataCredential_Int:
				dataCred[v.Data.GetKey()] = vv.Int
			case *pb.TxCredential_DataCredential_Str:
				dataCred[v.Data.GetKey()] = vv.Str
			}

		default:
			return txCredentials{}, fmt.Errorf("Unrecognized cred type %T", c.Cred)
		}
	}

	addrCred, err := NewAddrCredential(hash, ac)
	if err != nil {
		return txCredentials{}, err
	}

	return txCredentials{addrCred, dataCred}, nil
}

type dataCred map[string]interface{}

func (m dataCred) GetValue(k string) interface{} { return m[k] }

type AddrCredentialBuilder interface {
	AddSignature(sign *pb.Signature)
}

type DataCredentialsBuilder interface {
	AddString(string, string)
	AddInteger(string, int32)
	AddBytes(string, []byte)
}

type builder struct {
	e        error
	indexLen int
	cache    map[string]*pb.TxCredential_AddrCredentials
}

func newTxCredentialBuilder() *builder {
	return &builder{nil,
		fingerprintIndexLength,
		make(map[string]*pb.TxCredential_AddrCredentials),
	}
}

func (b *builder) AddSignature(sign *pb.Signature) {

	pk, err := crypto.PublicKeyFromSignature(sign)
	if err != nil {
		b.e = err
		return
	}

	ind := fmt.Sprintf("%.20X", pk.Digest())
	if _, ok := b.cache[ind]; ok {
		b.e = fmt.Errorf("Duplicated signature from same source")
		return
	}

	b.cache[ind] = &pb.TxCredential_AddrCredentials{
		Cred: &pb.TxCredential_AddrCredentials_User{
			User: &pb.TxCredential_UserCredential{Signature: sign}},
	}

}

func (b *builder) AddString(k string, v string) {
	b.cache[k] = &pb.TxCredential_AddrCredentials{
		Cred: &pb.TxCredential_AddrCredentials_Data{
			Data: &pb.TxCredential_DataCredential{
				Key:  k,
				Data: &pb.TxCredential_DataCredential_Str{Str: v},
			},
		},
	}
}

func (b *builder) AddInteger(k string, v int32) {
	b.cache[k] = &pb.TxCredential_AddrCredentials{
		Cred: &pb.TxCredential_AddrCredentials_Data{
			Data: &pb.TxCredential_DataCredential{
				Key:  k,
				Data: &pb.TxCredential_DataCredential_Int{Int: v},
			},
		},
	}
}

func (b *builder) AddBytes(k string, v []byte) {
	b.cache[k] = &pb.TxCredential_AddrCredentials{
		Cred: &pb.TxCredential_AddrCredentials_Data{
			Data: &pb.TxCredential_DataCredential{
				Key:  k,
				Data: &pb.TxCredential_DataCredential_Bts{Bts: v},
			},
		},
	}
}

func (b *builder) update(msg *pb.TxCredential) error {

	if b.e != nil {
		return b.e
	}

	msg.Addrc = make([]*pb.TxCredential_AddrCredentials, 0, len(b.cache))

	for _, v := range b.cache {
		msg.Addrc = append(msg.Addrc, v)
	}

	return nil
}
