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
	GetDataCredBuilder() DataCredentialsBuilder
	GenArguments() ([][]byte, error)
	GenArgumentsWithoutCred() ([][]byte, error)
	GenHash(string) []byte
}

type baseMaker struct {
	txArgs      [][]byte
	credBuilder *builder
}

func NewTxMaker(args [][]byte) *baseMaker {
	return &baseMaker{args, newTxCredentialBuilder()}
}

func (b *baseMaker) GenHash(method string) []byte {

	return genHash(b.txArgs[0], b.txArgs[1], method)
}

func (b *baseMaker) GetCredBuilder() AddrCredentialBuilder {

	return b.credBuilder
}

func (b *baseMaker) GetDataCredBuilder() DataCredentialsBuilder {

	return b.credBuilder
}

func (b *baseMaker) GenArguments() ([][]byte, error) {
	if b.credBuilder == nil {
		return nil, errors.New("No cred yet")
	}

	cred := &pb.TxCredential{}
	err := b.credBuilder.update(cred)
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
	method string
	nonce  []byte
	txHash []byte
}

func (b *baseBuilder) GetHash() []byte {
	if b.txHash == nil {
		b.txHash = b.GenHash(b.method)
	}
	return b.txHash
}

func (b *baseBuilder) GetNonce() []byte {
	return b.nonce
}

func (b *baseBuilder) SetHeader(h *pb.TxHeader) error {

	hh := msgToByte(h)
	if hh == nil {
		return errors.New("Encode header fail")
	}

	b.nonce = h.GetNonce()
	b.txHash = nil

	b.txArgs[0] = hh

	return nil
}

func (b *baseBuilder) SetMethod(m string) {
	b.txHash = nil
	b.method = m
}

func (b *baseBuilder) SetMessage(msg proto.Message) error {

	hm := msgToByte(msg)
	if hm == nil {
		return errors.New("No message binded")
	}
	b.txArgs[1] = hm
	b.txHash = nil
	return nil
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

func newTxBuilder() *baseBuilder {

	return &baseBuilder{baseMaker: baseMaker{make([][]byte, 2), newTxCredentialBuilder()}}
}

func NewTxBuilderWithTimeLock(ccname string, nonce []byte, t time.Time) (b *baseBuilder, err error) {

	if nonce == nil {
		nonce = GenerateNonce()
	}

	header := &pb.TxHeader{
		Base: &pb.TxBase{
			Network: DefaultNetworkName(),
			Ccname:  ccname,
		},
		ExpiredTs: &timestamp.Timestamp{
			Seconds: t.Unix(),
			Nanos:   0},
		Nonce: nonce,
		Flags: TxFlag_Timelock().U(),
	}

	b = newTxBuilder()
	err = b.SetHeader(header)

	return
}

func NewTxBuilder2(ccname string, nonce []byte) (b *baseBuilder, err error) {

	if nonce == nil {
		nonce = GenerateNonce()
	}

	header := &pb.TxHeader{
		Base: &pb.TxBase{
			Network: DefaultNetworkName(),
			Ccname:  ccname,
		},
		ExpiredTs: &timestamp.Timestamp{
			Seconds: time.Now().Unix() + int64(queryEffectInHour*3600),
			Nanos:   0},
		Nonce: nonce,
	}

	b = newTxBuilder()
	err = b.SetHeader(header)

	return
}

func NewTxBuilder(ccname string, nonce []byte, method string, msg proto.Message) (*baseBuilder, error) {

	b, err := NewTxBuilder2(ccname, nonce)
	if err != nil {
		return nil, err
	}

	err = b.SetMessage(msg)
	if err != nil {
		return nil, err
	}

	b.SetMethod(method)

	return b, nil
}
