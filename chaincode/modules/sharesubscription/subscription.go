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
	WeightBase = 1000000
)

func newContract(contract map[string]uint32, delePk *crypto.PublicKey) (*pb.Contract, error) {

	pcon := &pb.Contract{
		DelegatorPk: delePk.PBMessage(),
	}

	var totalweight uint64
	var usedweight uint32
	var keys = make([]string, 0, len(contract))

	for k, weight := range contract {
		totalweight = totalweight + uint64(weight)
		keys = append(keys, k)
	}

	//the data MUST be added in a deterministic manner
	sort.Strings(keys)

	for _, saddr := range keys {

		weight := contract[saddr]
		//turn the weight into a base of weightBase
		weight = uint32(uint64(weight) * uint64(WeightBase) / totalweight)
		usedweight = weight + usedweight

		pcon.Status = append(pcon.Status, &pb.Contract_MemberStatus{weight, nil, saddr})
	}

	if usedweight < uint32(WeightBase) {

		resident := uint32(WeightBase) - usedweight

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

	} else if usedweight > uint32(WeightBase) {
		panic("The arithmetic in golang may have ruined")
	}

	return pcon, nil
}

func hashContract(contract *pb.Contract, nonce []byte) (*tx.Address, error) {

	pk, err := crypto.PublicKeyFromPBMessage(contract.DelegatorPk)
	if err != nil {
		return nil, err
	}

	pkAddr, err := tx.NewAddress(pk)
	if err != nil {
		return nil, err
	}

	var maphash []byte
	for _, v := range contract.Status {
		maphashItem, err := utils.DoubleSHA256(
			bytes.Join([][]byte{
				[]byte(v.MemberID),
				[]byte(strconv.Itoa(int(v.Weight)))}, nil))

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

	hash, err := utils.DoubleSHA256(bytes.Join([][]byte{
		pkAddr.Hash, maphash, nonce}, nil))
	if err != nil {
		return nil, err
	}

	return tx.NewAddressFromHash(hash), nil
}

func (db *baseContractTx) New(contract map[string]uint32,
	delePk *crypto.PublicKey) ([]byte, error) {

	pcon, err := newContract(contract, delePk)
	if err != nil {
		return nil, err
	}

	conAddr, err := hashContract(pcon, db.nonce)
	if err != nil {
		return nil, err
	}

	//Ts is not the part of new contract!
	t, _ := db.stub.GetTxTime()
	pcon.ContractTs = utils.CreatePBTimestamp(t)

	err = db.Set(conAddr.ToString(), pcon)
	if err != nil {
		return nil, err
	}

	return conAddr.Hash, nil
}

func (db *baseContractTx) Query(addr []byte) (error, *pb.Contract) {
	conAddr := tx.NewAddressFromHash(addr)

	con := &pb.Contract{}

	err := db.Get(conAddr.ToString(), con)
	if err != nil {
		return err, nil
	}

	return nil, con
}

func (db *baseContractTx) QueryOne(conaddr []byte, addr []byte) (error, *pb.Contract) {

	err, data := db.Query(conaddr)
	if err != nil {
		return err, nil
	}

	saddr := tx.NewAddressFromHash(addr).ToString()

	for _, v := range data.Status {
		if v.MemberID == saddr {
			data.Status = []*pb.Contract_MemberStatus{v}
			return nil, data
		}
	}

	return errors.New("Not a member"), nil

}

func (db *baseContractTx) Redeem(conaddr []byte, addr []byte, amount *big.Int, redeemAddr []byte) ([]byte, error) {

	conAddrs := tx.NewAddressFromHash(conaddr).ToString()

	contract := &pb.Contract{}

	err := db.Get(conAddrs, contract)
	if err != nil {
		return nil, err
	}

	err, acc := db.token.Account(conaddr)
	if err != nil {
		return nil, err
	}

	member, ok := contract.Find(tx.NewAddressFromHash(addr).ToString())
	if !ok {
		return nil, errors.New("Not a member")
	}

	totalRedeem := toAmount(contract.TotalRedeem)
	curBalance := toAmount(acc.Balance)

	totalAsset := curBalance.Add(totalRedeem, curBalance)

	memberAsset := big.NewInt(int64(member.Weight))
	memberAsset = memberAsset.Mul(memberAsset, totalAsset).
		Div(memberAsset, big.NewInt(WeightBase))

	haveRedeem := toAmount(member.TotalRedeem)

	if memberAsset.Cmp(haveRedeem) <= 0 {
		return nil, errors.New("Could not redeem more")
	}

	canRedeem := memberAsset.Sub(memberAsset, haveRedeem)
	if (amount.IsUint64() && amount.Uint64() == 0) || canRedeem.Cmp(amount) < 0 {
		amount = canRedeem
	}

	member.TotalRedeem = haveRedeem.Add(haveRedeem, amount).Bytes()
	contract.TotalRedeem = totalRedeem.Add(totalRedeem, amount).Bytes()

	err = db.Set(conAddrs, contract)
	if err != nil {
		return nil, err
	}

	if redeemAddr == nil {
		redeemAddr = addr
	}

	return db.token.Transfer(conaddr, redeemAddr, amount)
}
