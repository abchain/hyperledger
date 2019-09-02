package abchainTx

import (
	"strings"
	"testing"

	abcrypto "hyperledger.abchain.org/core/crypto"
	"hyperledger.abchain.org/core/crypto/ecdsa"
	pb "hyperledger.abchain.org/protos"
)

// var privkey *abcrypto.PrivateKey

var privkey abcrypto.Signer

func txinit(t *testing.T) {

	var err error

	privkey, err = ecdsa.NewPrivatekey(ecdsa.DefaultCurveType)
	if err != nil {
		t.Fatal("Generate private key fail:", err)
	}
}

func signInBuilder(builder Builder, privk abcrypto.Signer) error {
	sig, err := privk.Sign(builder.GetHash())
	if err != nil {
		return err
	}

	builder.GetCredBuilder().AddSignature(sig)

	return nil
}

func TestTx(t *testing.T) {

	txinit(t)

	msg := &pb.TxMsgExample{Param1: []byte{'1', '9', '8', '4'}, Param2: 1984}

	builder, err := NewTxBuilder("gamepai", nil, "example", msg)

	if err != nil {
		t.Fatal("builder fail", err)
	}

	signInBuilder(builder, privkey)

	args, err := builder.GenArguments()

	if err != nil {
		t.Fatal("builder gen arg fail", err)
	}

	t.Log("General argument:", args)

	msgPending := &pb.TxMsgExample{}

	parser, err := ParseTx(msgPending, "example", args)

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
