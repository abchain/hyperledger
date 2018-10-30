package generaltoken

import (
	"hyperledger.abchain.org/chaincode/generaltoken/nonce"
	pb "hyperledger.abchain.org/chaincode/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/lib/util"
	"math/big"
)

type tokenTxNonce struct {
	Key     []byte
	Data    *pb.NonceData
	FromKey string
	ToKey   string
	From    *pb.AccountData
	To      *pb.AccountData
}

func (db *baseTokenTx) txNonce(txnonce []byte, from []byte, to []byte, amount *big.Int) (e error, r *tokenTxNonce) {

	abyte := amount.Bytes()

	r = &tokenTxNonce{
		Key: nonce.GeneralTokenNonceKey(txnonce, from, to, abyte),
		Data: &pb.NonceData{
			util.GetTxID(db.stub),
			abyte,
			nil, nil, util.GetTimeStamp(db.stub),
		},
	}

	if from != nil {
		r.FromKey = addrToKey(from)
		senderAcc := &pb.AccountData{}
		e = db.Get(r.FromKey, senderAcc)
		if e != nil {
			return
		}

		if senderAcc.Balance != nil {
			r.From = senderAcc
			r.Data.FromLast = senderAcc.LastFund
		}
	}

	if to != nil {
		r.ToKey = addrToKey(to)
		recvAcc := &pb.AccountData{}
		e = db.Get(r.ToKey, recvAcc)
		if e != nil {
			return
		}

		if recvAcc.Balance != nil {
			r.To = recvAcc
			r.Data.ToLast = recvAcc.LastFund
		}
	}

	return

}

func (db *baseTokenTx) Nonce(key []byte) (error, *pb.NonceData) {
	return db.tokenNonce.Nonce(key)
}

func (db *baseTokenTx) Add(r []byte, data *pb.NonceData) error {
	return db.tokenNonce.Add(r, data)
}
