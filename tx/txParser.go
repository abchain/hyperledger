package abchainTx

import (
	"errors"
	"github.com/abchain/fabric/core/chaincode/shim"
	"github.com/golang/protobuf/proto"
	pb "hyperledger.abchain.org/protos"
	"strings"
	"time"
)

type Parser interface {
	GetCCname() string
	GetNounce() []byte
	GetTxTime() time.Time
	GetMessage() proto.Message
	GetAddrCredential() AddrCredentials
}

type txParser struct {
	msg    proto.Message
	nonce  []byte
	ccname string
	txts   time.Time
	cred   AddrCredentials
}

func (t *txParser) GetCCname() string {
	return t.ccname
}

func (t *txParser) GetNounce() []byte {
	return t.nonce
}

func (t *txParser) GetTxTime() time.Time {
	return t.txts
}

func (t *txParser) GetMessage() proto.Message {
	return t.msg
}

func (t *txParser) GetAddrCredential() AddrCredentials {
	return t.cred
}

func parseBase(header proto.Message, msg proto.Message,
	stub shim.ChaincodeStubInterface, method string, args []string) (e error,
	cred AddrCredentials) {

	if len(args) < 2 {
		e = errors.New("No enough arguments")
		return
	}

	hh := fromArgument(args[0])

	if hh == nil {
		e = errors.New("Invalid header")
		return
	}

	e = proto.Unmarshal(hh, header)
	if e != nil {
		return
	}

	hm := fromArgument(args[1])

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

		hc := fromArgument(args[2])
		if hm == nil {
			e = errors.New("Invalid credential")
			return
		}

		credData := &pb.TxCredential{}

		e = proto.Unmarshal(hc, credData)
		if e != nil {
			return
		}

		cred, e = NewAddrCredential(hash, stub, credData.Addrc)

	} else {
		cred, e = NewAddrCredential(hash, stub, nil)
	}

	return
}

func ParseTx(msg proto.Message, stub shim.ChaincodeStubInterface, method string, args []string) (Parser, error) {

	header := &pb.TxHeader{}

	err, cred := parseBase(header, msg, stub, method, args)

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
		msg,
		header.Nonce,
		header.Base.Ccname,
		txTs,
		cred,
	}, nil

}
