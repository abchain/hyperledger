package abchainTx

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	pb "hyperledger.abchain.org/protos"
)

type TxMaker interface {
	GetCredBuilder() AddrCredentialBuilder
	GenArguments() ([][]byte, error)
	GenArgumentsWithoutCred() ([][]byte, error)
}

type baseMaker struct {
	txArgs      [][]byte
	credBuilder AddrCredentialBuilder
}

func NewTxMaker(args [][]byte) TxMaker {
	return &baseMaker{args, NewAddrCredentialBuilder()}
}

func (b *baseMaker) GetCredBuilder() AddrCredentialBuilder {

	return b.credBuilder
}

func (b *baseMaker) GenArguments() ([][]byte, error) {
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

	return append(arg, cr), nil
}

func (b *baseMaker) GenArgumentsWithoutCred() ([][]byte, error) {
	return b.txArgs, nil
}

type Builder interface {
	TxMaker
	GetNonce() []byte
	GetHash() []byte
}

type baseBuilder struct {
	baseMaker
	txHash []byte
	nonce  []byte
}

func (b *baseBuilder) GetHash() []byte {
	return b.txHash
}

func (b *baseBuilder) GetNonce() []byte {
	return b.nonce
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

	hm := msgToByte(msg)
	if hm == nil {
		return nil, errors.New("No message binded")
	}

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

	hh := msgToByte(header)
	if hh == nil {
		return nil, errors.New("Encode header fail")
	}

	b := &baseBuilder{
		baseMaker{[][]byte{hh, hm}, NewAddrCredentialBuilder()},
		genHash(hh, hm, method),
		nonce}

	return b, nil

}
