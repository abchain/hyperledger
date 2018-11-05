package abchainTx

import (
	"crypto/rand"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"hyperledger.abchain.org/core/crypto"
	pb "hyperledger.abchain.org/protos"
	"time"
)

type Builder interface {
	GetNonce() []byte
	GetHash() []byte
	GetCredBuilder() AddrCredentialBuilder
	Sign(*crypto.PrivateKey) error

	GenArguments() ([]string, error)
	GenArgumentsWithoutCred() ([]string, error)
}

type baseBuilder struct {
	tx
	txHash      []byte
	nonce       []byte
	credBuilder AddrCredentialBuilder
}

func (b *baseBuilder) GetCredBuilder() AddrCredentialBuilder {
	return b.credBuilder
}

func (b *baseBuilder) GetHash() []byte {
	return b.txHash
}

func (b *baseBuilder) GetNonce() []byte {
	return b.nonce
}

func (b *baseBuilder) Sign(privk *crypto.PrivateKey) error {

	if b.credBuilder == nil {
		b.credBuilder = NewAddrCredentialBuilder()
	}

	pk := privk.Public()

	sig, err := privk.SignwithThis(rand.Reader, b.txHash)
	if err != nil {
		return err
	}

	b.credBuilder.AddSignature(pk, &crypto.ECSignature{*sig})

	return nil
}

func (b *baseBuilder) GenArguments() ([]string, error) {
	if b.credBuilder == nil {
		return nil, errors.New("No cred yet")
	}

	cred := &pb.TxCredential{}
	err := b.credBuilder.Update(cred)
	if err != nil {
		return nil, err
	}

	cr := msgToByte(cred)
	if cr == nil {
		return nil, errors.New("Invalid cred data")
	}

	arg, err := b.GenArgumentsWithoutCred()

	if err != nil {
		return nil, err
	}

	return append(arg, toArgument(cr)), nil
}

func (b *baseBuilder) GenArgumentsWithoutCred() ([]string, error) {

	if b.msgObj == nil {
		return nil, errors.New("No message binded")
	}

	hh := msgToByte(b.header)
	hm := msgToByte(b.msgObj)

	if hh == nil || hm == nil {
		return nil, errors.New("Invalid message data")
	}

	return []string{toArgument(hh), toArgument(hm)}, nil
}

const (
	queryEffectInHour int = 1
)

func GenerateNonce() []byte {

	nonce := make([]byte, 20)
	_, err := rand.Read(nonce)
	if err != nil {
		//try different way to generate a nonce
		return []byte(time.Now().String())
	}

	return nonce
}

func NewTxBuilder(ccname string, nonce []byte, method string, msg proto.Message) (Builder, error) {

	if nonce == nil {
		nonce = GenerateNonce()
	}

	expTime := &timestamp.Timestamp{
		Seconds: time.Now().Unix() + int64(queryEffectInHour*3600),
		Nanos:   0}

	header := &pb.TxHeader{
		&pb.TxBase{
			DefaultNetworkName(),
			ccname,
			"",
		},
		expTime,
		nonce,
	}

	b := &baseBuilder{tx{header, msg}, nil, nonce, nil}
	b.txHash = b.GenHash(method)

	return b, nil

}
