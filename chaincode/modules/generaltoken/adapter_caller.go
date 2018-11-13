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

type FullGeneralCall struct {
	*GeneralCall
	nonce.TokenNonceTx
}

func NewFullGeneralCall(core txgen.TxCaller) *FullGeneralCall {

	return &FullGeneralCall{&GeneralCall{core}, &nonce.GeneralCall{core}}
}

const (
	Method_Init        = "TOKEN.INIT"
	Method_Transfer    = "TOKEN.TRANSFER"
	Method_Assign      = "TOKEN.ASSIGN"
	Method_TouchAddr   = "TOKEN.TOUCHADDR"
	Method_QueryToken  = "TOKEN.BALANCEQUERY"
	Method_QueryGlobal = "TOKEN.GLOBALQUERY"
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

	return nonce.GeneralTokenNonceKey(i.GetNonce(), from, to), nil
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

	return nonce.GeneralTokenNonceKey(i.GetNonce(), nil, to), nil
}

func (i *GeneralCall) TouchAddr(to []byte) error {
	return i.Invoke(Method_TouchAddr,
		&pb.QueryToken{Addr: txutil.NewAddressFromHash(to).PBMessage()})
}

func (i *GeneralCall) Init(amount *big.Int) error {

	msg := &pb.BaseToken{TotalTokens: amount.Bytes()}

	return i.Invoke(Method_Init, msg)
}

func (i *GeneralCall) Account(addr []byte) (error, *pb.AccountData_s) {
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

	out := &pb.AccountData_s{}
	out.LoadFromPB(d)
	return nil, out
}

func (i *GeneralCall) Global() (error, *pb.TokenGlobalData_s) {

	ret, err := i.Query(Method_QueryGlobal, &pb.QueryGlobal{})

	if err != nil {
		return err, nil
	}

	d := &pb.TokenGlobalData{}

	err = txgen.SyncQueryResult(d, ret)
	if err != nil {
		return err, nil
	}

	out := &pb.TokenGlobalData_s{}
	out.LoadFromPB(d)
	return nil, out

}
