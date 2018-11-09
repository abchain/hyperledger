package subscription

import (
	"bytes"
	"errors"
	pb "hyperledger.abchain.org/chaincode/modules/sharesubscription/protos"
	"hyperledger.abchain.org/core/crypto"
	tx "hyperledger.abchain.org/core/tx"
	"hyperledger.abchain.org/core/utils"
	"math/big"
	"sort"
	"strconv"
)

func toAmount(a []byte) *big.Int {

	if a == nil {
		return big.NewInt(0)
	}

	return big.NewInt(0).SetBytes(a)
}

const (
	WeightBase    = 1000000
	MaxContractor = 1024
)

func newContract(contract map[string]int32, delePk *crypto.PublicKey) (*pb.Contract_s, error) {

	if len(contract) > MaxContractor {
		return nil, errors.New("Too many contractors")
	}

	pcon := &pb.Contract_s{}
	pcon.TotalRedeem = big.NewInt(0)

	fgp, err := crypto.GetPublicKeyRootFingerprint(delePk.Key)
	if err != nil {
		return nil, err
	}
	pcon.DelegatorPkFingerPrint = fgp

	var totalweight int64
	var usedweight int32
	var keys = make([]string, 0, len(contract))

	for k, weight := range contract {
		if weight < 0 {
			return nil, errors.New("Minus weight is not allowed")
		}
		totalweight = totalweight + int64(weight)
		keys = append(keys, k)
	}

	//the data MUST be added in a deterministic manner
	sort.Strings(keys)

	for _, saddr := range keys {

		weight := contract[saddr]
		//turn the weight into a base of weightBase
		weight = int32(int64(weight) * int64(WeightBase) / totalweight)
		usedweight = weight + usedweight

		pcon.Status = append(pcon.Status, pb.Contract_MemberStatus_s{weight, big.NewInt(0), saddr})
	}

	if usedweight < int32(WeightBase) {

		resident := int32(WeightBase) - usedweight

		for resident > 0 {
			//simply add them into first N contractors ...
			for k, _ := range pcon.Status {

				pcon.Status[k].Weight++
				resident--

				if resident == 0 {
					break
				}
			}
		}

	} else if usedweight > int32(WeightBase) {
		panic("The arithmetic in golang may have ruined")
	}

	return pcon, nil
}

func hashContract(contract *pb.Contract_s, nonce []byte) (*tx.Address, error) {

	var maphash []byte
	for _, v := range contract.Status {
		maphashItem, err := utils.HMACSHA256(
			bytes.Join([][]byte{
				[]byte(v.MemberID),
				[]byte(strconv.Itoa(int(v.Weight)))}, nil),
			contract.DelegatorPkFingerPrint)

		if err != nil {
			return nil, err
		}

		if maphash == nil {
			maphash = maphashItem
		} else {
			//XOR bytes
			if len(maphash) != len(maphashItem) {
				return nil, errors.New("Wrong hash len")
			}

			for i := 0; i < len(maphash); i++ {
				maphash[i] = maphash[i] ^ maphashItem[i]
			}
		}
	}

	hash, err := utils.DoubleSHA256(bytes.Join([][]byte{maphash, nonce}, nil))
	if err != nil {
		return nil, err
	}

	return tx.NewAddressFromHash(hash), nil
}

func (cn *baseContractTx) New(contract map[string]int32,
	delePk *crypto.PublicKey) ([]byte, error) {

	pcon, err := newContract(contract, delePk)
	if err != nil {
		return nil, err
	}

	conAddr, err := hashContract(pcon, cn.nonce)
	if err != nil {
		return nil, err
	}

	t, _ := cn.Tx.GetTxTime()
	pcon.ContractTs = t

	err = cn.Storage.Set(conAddr.ToString(), pcon)
	if err != nil {
		return nil, err
	}

	return conAddr.Hash, nil
}

func (cn *baseContractTx) Query(addr []byte) (error, *pb.Contract_s) {
	conAddr := tx.NewAddressFromHash(addr)

	con := &pb.Contract_s{}

	err := cn.Storage.Get(conAddr.ToString(), con)
	if err != nil {
		return err, nil
	}

	return nil, con
}

func (cn *baseContractTx) QueryOne(conaddr []byte, addr []byte) (error, *pb.Contract_s) {

	err, data := cn.Query(conaddr)
	if err != nil {
		return err, nil
	}

	m, ok := data.Find(tx.NewAddressFromHash(addr).ToString())
	if !ok {
		return errors.New("Not a member"), nil
	}

	data.Status = []pb.Contract_MemberStatus_s{m}
	return nil, data
}

func (cn *baseContractTx) Redeem(conaddr []byte, addr []byte, amount *big.Int, redeemAddr []byte) ([]byte, error) {

	conAddrs := tx.NewAddressFromHash(conaddr).ToString()

	contract := &pb.Contract_s{}

	err := cn.Storage.Get(conAddrs, contract)
	if err != nil {
		return nil, err
	}

	err, acc := cn.token.Account(conaddr)
	if err != nil {
		return nil, err
	}

	member, mindex := contract.FindAndAccess(tx.NewAddressFromHash(addr).ToString())
	if mindex < 0 {
		return nil, errors.New("Not a member")
	}

	totalAsset := big.NewInt(0).Add(contract.TotalRedeem, acc.Balance)

	memberAsset := big.NewInt(int64(member.Weight))
	memberAsset = memberAsset.Mul(memberAsset, totalAsset).
		Div(memberAsset, big.NewInt(WeightBase))

	if memberAsset.Cmp(member.TotalRedeem) <= 0 {
		return nil, errors.New("Could not redeem more")
	}

	canRedeem := big.NewInt(0).Sub(memberAsset, member.TotalRedeem)
	if (amount.Int64() == 0) || canRedeem.Cmp(amount) < 0 {
		amount = canRedeem
	}

	contract.TotalRedeem = big.NewInt(0).Add(contract.TotalRedeem, amount)
	contract.Status[mindex].TotalRedeem = big.NewInt(0).Add(member.TotalRedeem, amount)

	err = cn.Storage.Set(conAddrs, contract)
	if err != nil {
		return nil, err
	}

	if redeemAddr == nil {
		redeemAddr = addr
	}

	return cn.token.Transfer(conaddr, redeemAddr, amount)
}
