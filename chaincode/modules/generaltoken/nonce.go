package generaltoken

import (
	"hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"hyperledger.abchain.org/utils"
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
	t, _ := db.stub.GetTxTime()

	r = &tokenTxNonce{
		Key: nonce.GeneralTokenNonceKey(txnonce, from, to, abyte),
		Data: &pb.NonceData{
			db.stub.GetTxID(),
			abyte,
			nil, nil, utils.CreatePBTimestamp(t),
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
