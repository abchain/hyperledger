package generaltoken

import (
	"hyperledger.abchain.org/chaincode/generaltoken/nonce"
	pb "hyperledger.abchain.org/chaincode/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	txutil "hyperledger.abchain.org/tx"
	"math/big"
)

type GeneralCall struct {
	*txgen.TxGenerator
}

const (
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

	_, err := i.Invoke(Method_Transfer, msg)
	if err != nil {
		return nil, err
	}

	return nonce.GeneralTokenNonceKey(i.GetBuilder().GetNonce(), from, to, amount.Bytes()), nil
}

func (i *GeneralCall) Assign(to []byte, amount *big.Int) ([]byte, error) {

	msg := &pb.SimpleFund{
		amount.Bytes(),
		txutil.NewAddressFromHash(to).PBMessage(),
		nil,
	}

	_, err := i.Invoke(Method_Assign, msg)
	if err != nil {
		return nil, err
	}

	return nonce.GeneralTokenNonceKey(i.GetBuilder().GetNonce(), nil, to, amount.Bytes()), nil
}

func (i *GeneralCall) Account(addr []byte) (error, *pb.AccountData) {

	a := txutil.NewAddressFromHash(addr)
	ret, err := i.Query(Method_QueryToken, &pb.QueryToken{pb.QueryToken_ENCODED, a.PBMessage()})

	if err != nil {
		return err, nil
	}

	if ret == nil {
		return nil, nil
	}

	d := &pb.AccountData{}

	err = rpc.DecodeRPCResult(d, ret)
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

	if ret == nil {
		return nil, nil
	}

	d := &pb.NonceData{}

	err = rpc.DecodeRPCResult(d, ret)
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

	if ret == nil {
		return nil, nil
	}

	d := &pb.TokenGlobalData{}

	err = rpc.DecodeRPCResult(d, ret)
	if err != nil {
		return err, nil
	}

	return nil, d

}
