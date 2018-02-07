package abchainTx

import (
	abcrypto "hyperledger.abchain.org/crypto"
	pb "hyperledger.abchain.org/protos"
	"strings"
	"testing"
)

var privkey *abcrypto.PrivateKey

func txinit(t *testing.T) {

	var err error

	privkey, err = abcrypto.NewPrivatekey(abcrypto.DefaultCurveType)
	if err != nil {
		t.Fatal("Generate private key fail:", err)
	}
}

func TestTx(t *testing.T) {

	txinit(t)

	msg := &pb.TxMsgExample{[]byte{'1', '9', '8', '4'}, 1984}

	builder, err := NewTxBuilder("gamepai", nil, "example", msg)

	if err != nil {
		t.Fatal("builder fail", err)
	}

	builder.Sign(privkey)

	args, err := builder.GenArguments()

	if err != nil {
		t.Fatal("builder gen arg fail", err)
	}

	t.Log("General argument:", args)

	msgPending := &pb.TxMsgExample{}

	parser, err := ParseTx(msgPending, nil, "example", args)

	if err != nil {
		t.Fatal("Parse arg fail", err)
	}

	if strings.Compare(parser.GetCCname(), "gamepai") != 0 {
		t.Fatal("Wrong ccname")
	}

	if msgPending.Param2 != msg.Param2 {
		t.Fatal("Wrong message")
	}
}
