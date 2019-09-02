package abchainTx

import (
	"errors"
	"github.com/golang/protobuf/proto"
	pb "hyperledger.abchain.org/protos"
	"strings"
	"time"
)

type Flags uint32

//indicate the expried time in header is obliged to be considered
func TxFlag_Timelock() Flags {
	return Flags(1)
}

func (f Flags) U() uint32 {
	return uint32(f)
}

func (f Flags) And(fanother Flags) Flags {
	return Flags(uint32(f) | uint32(fanother))
}

func (f Flags) IsTimeLock() bool {
	return (f & TxFlag_Timelock()) != 0
}

type Parser interface {
	GetCCname() string
	GetNonce() []byte
	GetFlags() Flags
	GetTxTime() time.Time
	GetAddrCredential() AddrCredentials
	GetDataCredential() DataCredentials
	GetMessage() proto.Message
	UpdateMsg(proto.Message)
}

type txParser struct {
	nonce  []byte
	ccname string
	flags  uint32
	txts   time.Time
	cred   txCredentials
	msg    proto.Message
}

func (t *txParser) GetCCname() string {
	return t.ccname
}

func (t *txParser) GetNonce() []byte {
	return t.nonce
}

func (t *txParser) GetFlags() Flags {
	return Flags(t.flags)
}

func (t *txParser) GetTxTime() time.Time {
	return t.txts
}

func (t *txParser) GetAddrCredential() AddrCredentials {
	return t.cred
}

func (t *txParser) GetDataCredential() DataCredentials {
	return t.cred
}

func (t *txParser) GetMessage() proto.Message {
	return t.msg
}

func (t *txParser) UpdateMsg(m proto.Message) {
	t.msg = m
}

func parseBase(header proto.Message, msg proto.Message,
	method string, args [][]byte) (e error, cred txCredentials) {

	if len(args) < 2 {
		e = errors.New("No enough arguments")
		return
	}

	hh := args[0]

	if hh == nil {
		e = errors.New("Invalid header")
		return
	}

	e = proto.Unmarshal(hh, header)
	if e != nil {
		return
	}

	hm := args[1]

	if hm == nil {
		e = errors.New("Invalid msg")
		return
	}

	e = proto.Unmarshal(hm, msg)
	if e != nil {
		return
	}

	hash := genHash(hh, hm, method)
	if hash == nil {
		e = errors.New("Hashing raw tx fail")
		return
	}

	if len(args) >= 3 {

		hc := args[2]
		if hm == nil {
			e = errors.New("Invalid credential")
			return
		}

		credData := &pb.TxCredential{}

		e = proto.Unmarshal(hc, credData)
		if e != nil {
			return
		}

		cred, e = newTxCredential(hash, credData.Addrc)

	} else {
		cred, e = newTxCredential(hash, nil)
	}

	return
}

func ParseTx(msg proto.Message, method string, args [][]byte) (Parser, error) {

	header := &pb.TxHeader{}

	err, cred := parseBase(header, msg, method, args)

	if err != nil {
		return nil, err
	}

	if strings.Compare(header.Base.Network, DefaultNetworkName()) != 0 {
		return nil, errors.New("Unmatch network")
	}

	var txTs time.Time
	if header.ExpiredTs != nil {
		txTs = time.Unix(header.ExpiredTs.Seconds, int64(header.ExpiredTs.Nanos))
	}

	return &txParser{
		header.Nonce,
		header.Base.Ccname,
		header.Flags,
		txTs,
		cred,
		msg,
	}, nil

}
