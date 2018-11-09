package nonce

import (
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"math/big"
)

type GeneralCall struct {
	txgen.TxCaller
}

const (
	Method_Add   = "NONCE.ADD"
	Method_Query = "NONCE.QUERY"
)

func (i *GeneralCall) Add(key []byte, amount *big.Int, from *pb.FuncRecord, to *pb.FuncRecord) error {

	//use noncedata to carry information
	return i.Invoke(Method_Add, &pb.NonceData{
		Amount:   amount.Bytes(),
		FromLast: from,
		ToLast:   to,
		Other:    &pb.NonceData_Noncekey{Noncekey: key},
	})
}

func (i *GeneralCall) Nonce(key []byte) (error, *pb.NonceData_s) {

	ret, err := i.Query(Method_Query, &pb.QueryTransfer{key})

	if err != nil {
		return err, nil
	}

	d := &pb.NonceData{}

	err = txgen.SyncQueryResult(d, ret)
	if err != nil {
		return err, nil
	}

	out := &pb.NonceData_s{}
	out.LoadFromPB(d)
	return nil, out
}
