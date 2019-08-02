package generaltoken

import (
	"errors"
	"math/big"

	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
)

const (
	deployName = ":deploy"
)

func (token *baseTokenTx) Assign(to []byte, amount *big.Int) (pb.NonceKey, error) {
	global := &pb.TokenGlobalData_s{}
	// global := &model.TokenGlobalData{}
	err := token.Storage.Get(deployName, global)
	if err != nil {
		return nil, err
	}

	if global.TotalTokens == nil {
		return nil, errors.New("Not deploy yet")
	}

	err, ret := token.txNonce(token.nonce, nil, to, amount)

	if err != nil {
		return nil, err
	}

	var fromLast, toLast *pb.FuncRecord
	if ret.From != nil {
		fromLast = ret.From.LastFund.ToPB()
	}
	if ret.To != nil {
		toLast = ret.To.LastFund.ToPB()
	}

	err = token.tokenNonce.Add(ret.Key, amount, fromLast, toLast)
	if err != nil {
		return nil, err
	}

	if global.UnassignedTokens.Cmp(amount) < 0 {
		return nil, errors.New("Total amount is not enough for assigment")
	}

	global.UnassignedTokens = big.NewInt(0).Sub(global.UnassignedTokens, amount)

	if ret.To == nil {
		ret.To = &pb.AccountData_s{}
		ret.To.Balance = amount
	} else {
		ret.To.Balance = big.NewInt(0).Add(ret.To.Balance, amount)
	}

	ret.To.LastFund = pb.FuncRecord_s{ret.Key, false}

	// err = token.Storage.Set(ret.ToKey, &pb.AccountData{
	// 	balance.Add(balance, amount).Bytes(),
	// 	&pb.FuncRecord{ret.Key, false},
	// })

	err = token.Storage.Set(ret.ToKey, ret.To)
	if err != nil {
		return nil, err
	}

	err = token.Storage.Set(deployName, global)
	if err != nil {
		return nil, err
	}
	return ret.Key, nil
}

func (token *baseTokenTx) Global() (error, *pb.TokenGlobalData_s) {

	global := &pb.TokenGlobalData_s{}

	err := token.Storage.Get(deployName, global)

	if err != nil {
		return err, nil
	}

	if global.TotalTokens == nil {
		return errors.New("token is null , you need init token first"), nil
	}

	return nil, global
}
