package generaltoken

import (
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	txutil "hyperledger.abchain.org/core/tx"
	"math/big"
)

type GeneralCall struct {
	txgen.TxCaller
}

const (
	Method_Init        = "TOKEN.INIT"
	Method_Transfer    = "TOKEN.TRANSFER"
	Method_Assign      = "TOKEN.ASSIGN"
	Method_QueryToken  = "TOKEN.BALANCEQUERY"
	Method_QueryGlobal = "TOKEN.GLOBALQUERY"
	Method_QueryTrans  = "TOKEN.TRANSFERQUERY"
)

func (i *GeneralCall) Transfer(from []byte, to []byte, amount *big.Int) ([]byte, error) {

	msg := &pb.SimpleFund{
		amount.Bytes(),
		txutil.NewAddressFromHash(to).PBMessage(),
		txutil.NewAddressFromHash(from).PBMessage(),
	}

	err := i.Invoke(Method_Transfer, msg)
	if err != nil {
		return nil, err
	}

	return nonce.GeneralTokenNonceKey(i.GetNonce(), from, to, amount.Bytes()), nil
}

func (i *GeneralCall) Assign(to []byte, amount *big.Int) ([]byte, error) {

	msg := &pb.SimpleFund{
		amount.Bytes(),
		txutil.NewAddressFromHash(to).PBMessage(),
		nil,
	}

	err := i.Invoke(Method_Assign, msg)
	if err != nil {
		return nil, err
	}

	return nonce.GeneralTokenNonceKey(i.GetNonce(), nil, to, amount.Bytes()), nil
}

func (i *GeneralCall) Init(amount *big.Int) error {

	msg := &pb.BaseToken{TotalTokens: amount.Bytes()}

	return i.Invoke(Method_Init, msg)
}

func (i *GeneralCall) Account(addr []byte) (error, *pb.AccountData) {

	a := txutil.NewAddressFromHash(addr)
	ret, err := i.Query(Method_QueryToken, &pb.QueryToken{pb.QueryToken_ENCODED, a.PBMessage()})

	if err != nil {
		return err, nil
	}

	d := &pb.AccountData{}

	err = txgen.SyncQueryResult(d, ret)
	if err != nil {
		return err, nil
	}

	return nil, d
}

func (i *GeneralCall) Nonce(key []byte) (error, *pb.NonceData) {

	ret, err := i.Query(Method_QueryTrans, &pb.QueryTransfer{key})

	if err != nil {
		return err, nil
	}

	d := &pb.NonceData{}

	err = txgen.SyncQueryResult(d, ret)
	if err != nil {
		return err, nil
	}

	return nil, d
}

func (i *GeneralCall) Global() (error, *pb.TokenGlobalData) {

	ret, err := i.Query(Method_QueryGlobal, &pb.QueryGlobal{})

	if err != nil {
		return err, nil
	}

	d := &pb.TokenGlobalData{}

	err = txgen.SyncQueryResult(d, ret)
	if err != nil {
		return err, nil
	}

	return nil, d

}
