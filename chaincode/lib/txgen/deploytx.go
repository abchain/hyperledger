package tx

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txutil "hyperledger.abchain.org/tx"
)

type DeployTxCall struct {
	*TxGenerator
}

func (i *GeneralCall) BuildDeployArg(from []byte, to []byte, amount *big.Int) ([]byte, error) {

	msg := &pb.SimpleFund{
		amount.Bytes(),
		txutil.NewAddressFromHash(to).PBMessage(),
		txutil.NewAddressFromHash(from).PBMessage(),
	}

	_, err := i.Invoke(Method_Transfer, msg)
	if err != nil {
		return nil, err
	}

	return nonce.GeneralTokenNonceKey(i.GetBuilder().GetNonce(), from, to, amount.Bytes()), nil
}
