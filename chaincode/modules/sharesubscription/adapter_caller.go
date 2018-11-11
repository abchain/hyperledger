package subscription

import (
	"errors"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
	pb "hyperledger.abchain.org/chaincode/modules/sharesubscription/protos"
	txutil "hyperledger.abchain.org/core/tx"
	txpb "hyperledger.abchain.org/protos"
	"math/big"
)

type GeneralCall struct {
	*txgen.TxGenerator
}

const (
	Method_NewContract = "CONTRACT.SUBSCRIPTION.NEW"
	Method_Redeem      = "CONTRACT.SUBSCRIPTION.REDEEM"
	Method_Query       = "CONTRACT.SUBSCRIPTION.QUERY"
	Method_MemberQuery = "CONTRACT.SUBSCRIPTION.QUERYONE"
)

func (i *GeneralCall) New(contract map[string]int32, addr []byte) ([]byte, error) {

	if len(contract) == 0 {
		return nil, errors.New("Empty contract")
	}

	contractTx := make([]*pb.RegContract_Member, 0, len(contract))

	for a, weight := range contract {
		addr, err := txutil.NewAddressFromString(a)
		if err != nil {
			return nil, err
		}

		contractTx = append(contractTx, &pb.RegContract_Member{addr.PBMessage(), weight})
	}

	msg := &pb.RegContract{&txpb.TxAddr{Hash: addr}, contractTx}
	err := i.Invoke(Method_NewContract, msg)
	if err != nil {
		return nil, err
	}

	//gen the contract addr
	data, err := newContract(contract, addr)
	if err != nil {
		return nil, err
	}

	conaddr, err := hashContract(data, i.GetNonce())
	if err != nil {
		return nil, err
	}

	return conaddr.Hash, err
}

func (i *GeneralCall) Redeem(conaddr []byte, amount *big.Int, redeemAddrs [][]byte) (*pb.RedeemResponse, error) {

	msg := &pb.RedeemContract{
		txutil.NewAddressFromHash(conaddr).PBMessage(),
		amount.Bytes(),
		nil,
	}

	ret := &pb.RedeemResponse{}
	for _, addr := range redeemAddrs {
		msg.Redeems = append(msg.Redeems, txutil.NewAddressFromHash(addr).PBMessage())
		ret.Nonces = append(ret.Nonces, nonce.GeneralTokenNonceKey(i.GetNonce(), conaddr, addr))
	}

	err := i.Invoke(Method_Redeem, msg)
	if err != nil {
		return nil, err
	}

	return ret, nil

}

func (i *GeneralCall) Query(addr []byte) (error, *pb.Contract_s) {

	msg := &pb.QueryContract{
		txutil.NewAddressFromHash(addr).PBMessage(),
		nil,
	}

	data, err := i.TxGenerator.Query(Method_Query, msg)
	if err != nil {
		return err, nil
	}

	d := &pb.Contract{}
	err = txgen.SyncQueryResult(d, data)
	if err != nil {
		return err, nil
	}

	ret := new(pb.Contract_s)
	ret.LoadFromPB(d)
	return nil, ret
}

func (i *GeneralCall) QueryOne(conaddr []byte, addr []byte) (error, *pb.Contract_s) {
	msg := &pb.QueryContract{
		txutil.NewAddressFromHash(conaddr).PBMessage(),
		txutil.NewAddressFromHash(addr).PBMessage(),
	}

	data, err := i.TxGenerator.Query(Method_MemberQuery, msg)
	if err != nil {
		return err, nil
	}

	d := &pb.Contract{}
	err = txgen.SyncQueryResult(d, data)
	if err != nil {
		return err, nil
	}

	ret := new(pb.Contract_s)
	ret.LoadFromPB(d)
	return nil, ret
}
