package generaltoken

import (
	"hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"math/big"
)

type tokenTxNonce struct {
	Key     []byte
	FromKey string
	ToKey   string
	From    *pb.AccountData_s
	To      *pb.AccountData_s
}

func (token *baseTokenTx) txNonce(txnonce []byte, from []byte, to []byte, amount *big.Int) (e error, r *tokenTxNonce) {

	r = &tokenTxNonce{
		Key: nonce.GeneralTokenNonceKey(txnonce, from, to, amount.Bytes()),
	}

	if from != nil {
		r.FromKey = addrToKey(from)
		senderAcc := &pb.AccountData_s{}
		e = token.Storage.Get(r.FromKey, senderAcc)
		if e != nil {
			return
		}

		if senderAcc.Balance != nil {
			r.From = senderAcc
		}
	}

	if to != nil {
		r.ToKey = addrToKey(to)
		recvAcc := &pb.AccountData_s{}
		e = token.Storage.Get(r.ToKey, recvAcc)
		if e != nil {
			return
		}

		if recvAcc.Balance != nil {
			r.To = recvAcc
		}
	}
	return

}

func (token *baseTokenTx) Nonce(key []byte) (error, *pb.NonceData) {
	return token.tokenNonce.Nonce(key)
}

func (token *baseTokenTx) Add(key []byte, amount *big.Int, from *pb.FuncRecord, to *pb.FuncRecord) error {
	return token.tokenNonce.Add(key, amount, from, to)
}
