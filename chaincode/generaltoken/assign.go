package generaltoken

import (
	"errors"
	pb "hyperledger.abchain.org/chaincode/generaltoken/protos"
	"math/big"
)

func (db *baseTokenTx) Assign(to []byte, amount *big.Int) ([]byte, error) {

	global := &pb.TokenGlobalData{}

	err := db.Get(deployName, global)
	if err != nil {
		return nil, err
	}

	if global.TotalTokens == nil {
		return nil, errors.New("Not deploy yet")
	}

	err, ret := db.txNonce(db.nonce, nil, to, amount)

	if err != nil {
		return nil, err
	}

	if db.tokenNonce != nil {
		err = db.tokenNonce.Add(ret.Key, ret.Data)
		if err != nil {
			return nil, err
		}
	}

	assigment := toAmount(global.UnassignedTokens)
	if assigment.Cmp(amount) < 0 {
		return nil, errors.New("Total amount is not enough for assigment")
	}

	assigment = assigment.Sub(assigment, amount)

	if ret.To == nil {
		ret.To = &pb.AccountData{}
	}

	balance := toAmount(ret.To.Balance)

	err = db.Set(ret.ToKey, &pb.AccountData{
		balance.Add(balance, amount).Bytes(),
		&pb.FuncRecord{ret.Key, false},
	})

	if err != nil {
		return nil, err
	}

	global.UnassignedTokens = assigment.Bytes()
	err = db.Set(deployName, global)
	if err != nil {
		return nil, err
	}

	return ret.Key, nil
}

func (db *baseTokenTx) Global() (error, *pb.TokenGlobalData) {

	global := &pb.TokenGlobalData{}

	err := db.Get(deployName, global)

	if err != nil {
		return err, nil
	}

	return nil, global
}
