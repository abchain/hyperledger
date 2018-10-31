package generaltoken

import (
	"errors"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"math/big"
)

func (db *baseTokenTx) Transfer(from []byte, to []byte, amount *big.Int) ([]byte, error) {

	err, ret := db.txNonce(db.nonce, from, to, amount)

	if err != nil {
		return nil, err
	}

	if ret.From == nil {
		return nil, errors.New("Fund sender is not exist")
	}

	if db.tokenNonce != nil {
		err = db.Add(ret.Key, ret.Data)
		if err != nil {
			return nil, err
		}
	}

	if ret.To == nil {
		ret.To = &pb.AccountData{}
	}

	senderBalance := toAmount(ret.From.Balance)
	recvBalance := toAmount(ret.To.Balance)

	if senderBalance.Cmp(amount) < 0 {
		return nil, errors.New("No enough balance")
	}

	senderBalance = senderBalance.Sub(senderBalance, amount)

	if senderBalance.Sign() < 0 {
		return nil, errors.New("Wrong balance!")
	}
	recvBalance = recvBalance.Add(recvBalance, amount)

	err = db.Set(ret.FromKey, &pb.AccountData{
		senderBalance.Bytes(),
		&pb.FuncRecord{ret.Key, true},
	})
	if err != nil {
		return nil, err
	}

	err = db.Set(ret.ToKey, &pb.AccountData{
		recvBalance.Bytes(),
		&pb.FuncRecord{ret.Key, false},
	})
	if err != nil {
		return nil, err
	}

	return ret.Key, nil

}

func (db *baseTokenTx) Account(addr []byte) (error, *pb.AccountData) {

	acc := &pb.AccountData{}
	err := db.Get(addrToKey(addr), acc)
	if err != nil {
		return err, nil
	}

	return nil, acc

}
