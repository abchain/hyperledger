package subscription

import (
	"encoding/base64"
	"errors"
	"hyperledger.abchain.org/chaincode/lib/runtime"
	pb "hyperledger.abchain.org/chaincode/modules/sharesubscription/protos"
	tx "hyperledger.abchain.org/core/tx"
	"hyperledger.abchain.org/core/utils"
	"math/big"
	"sort"
)

func addrToKey(h []byte) string {
	return base64.RawURLEncoding.EncodeToString(h)
}

func toAmount(a []byte) *big.Int {

	if a == nil {
		return big.NewInt(0)
	}

	return big.NewInt(0).SetBytes(a)
}

const (
	maxContractor = 128
)

func newContract(addrs [][]byte, ratios []int) (*pb.Contract_s, error) {

	if len(addrs) > maxContractor {
		return nil, errors.New("Too many contractors")
	} else if len(addrs) != len(ratios) {
		return nil, errors.New("Wrong arguments")
	}

	pcon := &pb.Contract_s{}
	pcon.TotalRedeem = big.NewInt(0)

	for i, ratio := range ratios {
		if ratio < 0 {
			return nil, errors.New("Minus weight is not allowed")
		}

		pcon.Status = append(pcon.Status,
			pb.Contract_MemberStatus_s{int32(ratio), big.NewInt(0), addrs[i]})
		pcon.TotalWeight += int64(ratio)
	}

	//the data MUST be added in a deterministic manner
	sort.Sort(pcon.Sorter())

	return pcon, nil
}

func hashContract(contract *pb.Contract_s, nonce []byte) ([]byte, error) {

	bts, err := runtime.SeralizeObject(contract)
	if err != nil {
		return nil, err
	}

	hash, err := utils.HMACSHA256(bts, nonce)
	if err != nil {
		return nil, err
	}

	if len(hash) > tx.ADDRESS_HASH_LEN {
		hash = hash[:tx.ADDRESS_HASH_LEN]
	}

	return hash, nil

}

func (cn *baseContractTx) New_C(addrs [][]byte, ratios []int) ([]byte, error) {

	pcon, err := newContract(addrs, ratios)
	if err != nil {
		return nil, err
	}

	hash, err := hashContract(pcon, cn.nonce)
	if err != nil {
		return nil, err
	}

	addrhash, err := cn.addrutil.NormalizeAddress(hash)
	if err != nil {
		return nil, err
	}

	t, _ := cn.Tx.GetTxTime()
	pcon.ContractTs = t

	err = cn.Storage.Set(addrToKey(addrhash), pcon)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func (cn *baseContractTx) Query_C(addr []byte) (error, *pb.Contract_s) {
	con := &pb.Contract_s{}

	err := cn.Storage.Get(addrToKey(addr), con)
	if err != nil {
		return err, nil
	} else if con.TotalRedeem == nil {
		return errors.New("Contract is not exist"), nil
	}

	return nil, con
}

func (cn *baseContractTx) QueryOne_C(conaddr, addr []byte) (error, *pb.Contract_s) {

	err, data := cn.Query_C(conaddr)
	if err != nil {
		return err, nil
	}

	m, ok := data.Find(addr)
	if !ok {
		return errors.New("Not a member"), nil
	}

	data.Status = []pb.Contract_MemberStatus_s{m}
	return nil, data
}

func (cn *baseContractTx) Redeem_C(conaddr []byte, amount *big.Int, redeemAddrs [][]byte) (*pb.RedeemResponse, error) {

	if len(redeemAddrs) == 0 {
		return nil, errors.New("No redeem addrs")
	}

	err, contract := cn.Query_C(conaddr)
	if err != nil {
		return nil, err
	}

	err, acc := cn.token.Account(conaddr)
	if err != nil {
		return nil, err
	}

	totalAsset := big.NewInt(0).Add(contract.TotalRedeem, acc.Balance)
	totalWeight := big.NewInt(contract.TotalWeight)

	ret := &pb.RedeemResponse{}
	for _, addr := range redeemAddrs {
		member, mindex := contract.FindAndAccess(addr)
		if mindex < 0 {
			continue
		}

		memberAsset := big.NewInt(int64(member.Weight))
		memberAsset = memberAsset.Mul(memberAsset, totalAsset).Div(memberAsset, totalWeight)
		if memberAsset.Cmp(member.TotalRedeem) <= 0 {
			continue
		}

		canRedeem := big.NewInt(0).Sub(memberAsset, member.TotalRedeem)
		if (amount.Int64() != 0) && canRedeem.Cmp(amount) > 0 {
			canRedeem = amount
		}

		if nc, err := cn.token.Transfer(conaddr, addr, canRedeem); err == nil {
			ret.Nonces = append(ret.Nonces, nc)
			contract.TotalRedeem = contract.TotalRedeem.Add(contract.TotalRedeem, canRedeem)
			contract.Status[mindex].TotalRedeem = big.NewInt(0).Add(member.TotalRedeem, canRedeem)
		}
	}

	err = cn.Storage.Set(addrToKey(conaddr), contract)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
