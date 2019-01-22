package generaltoken

import (
	"errors"
	"math/big"

	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
)

func (token *baseTokenTx) Transfer(from []byte, to []byte, amount *big.Int) ([]byte, error) {

	err, ret := token.txNonce(token.nonce, from, to, amount)

	if err != nil {
		return nil, err
	}

	if ret.From == nil {
		return nil, errors.New("Fund sender is not exist")
	}

	if ret.FromKey == ret.ToKey {
		return nil, errors.New("can not transfer asset to the same address")
	}

	var toLast *pb.FuncRecord
	if ret.To != nil {
		toLast = ret.To.LastFund.ToPB()
	}

	err = token.tokenNonce.Add(ret.Key, amount, ret.From.LastFund.ToPB(), toLast)
	if err != nil {
		return nil, err
	}

	if ret.To == nil {
		ret.To = &pb.AccountData_s{}
		ret.To.Balance = big.NewInt(0)
	}

	if ret.From.Balance.Cmp(amount) < 0 {
		return nil, errors.New("No enough balance")
	}

	ret.From.Balance = big.NewInt(0).Sub(ret.From.Balance, amount)
	ret.From.LastFund = pb.FuncRecord_s{ret.Key, true}

	if ret.From.Balance.Sign() < 0 {
		return nil, errors.New("Wrong balance!")
	}
	ret.To.Balance = big.NewInt(0).Add(ret.To.Balance, amount)
	ret.To.LastFund = pb.FuncRecord_s{ret.Key, false}

	err = token.Storage.Set(ret.FromKey, ret.From)
	if err != nil {
		return nil, err
	}

	err = token.Storage.Set(ret.ToKey, ret.To)
	if err != nil {
		return nil, err
	}

	return ret.Key, nil

}

func (token *baseTokenTx) Account(addr []byte) (error, *pb.AccountData_s) {

	acc := &pb.AccountData_s{}
	err := token.Storage.Get(addrToKey(addr), acc)
	if err != nil {
		return err, nil
	} else if acc.Balance == nil {
		acc.Balance = big.NewInt(0)
	}

	return nil, acc

}

func (token *baseTokenTx) TouchAddr(addr []byte) error {

	key := addrToKey(addr)
	acc := &pb.AccountData_s{}
	err := token.Storage.Get(key, acc)
	if err != nil {
		return err
	} else if acc.Balance == nil {
		acc.Balance = big.NewInt(0)
		return token.Storage.Set(key, acc)
	}

	return nil
}

func (token *baseTokenTx) Init(total *big.Int) error {
	deploy := &pb.TokenGlobalData_s{}
	err := token.Storage.Get(deployName, deploy)

	if deploy.TotalTokens != nil {
		return errors.New("Can not re-deploy existed data")
	}

	deploy.TotalTokens = total
	deploy.UnassignedTokens = total
	err = token.Storage.Set(deployName, deploy)
	if err != nil {
		return err
	}

	return nil
}
