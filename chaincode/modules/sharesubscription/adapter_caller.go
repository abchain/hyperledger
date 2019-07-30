package subscription

import (
	"errors"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"hyperledger.abchain.org/chaincode/modules/generaltoken/nonce"
	pb "hyperledger.abchain.org/chaincode/modules/sharesubscription/protos"
	txutil "hyperledger.abchain.org/core/tx"
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

func (i *GeneralCall) new(addrs [][]byte, ratios []int) ([]byte, *pb.RegContract, error) {

	if len(addrs) != len(ratios) {
		return nil, nil, errors.New("Wrong argument")
	}

	contractTx := make([]*pb.RegContract_Member, 0, len(addrs))

	for i, a := range addrs {

		contractTx = append(contractTx,
			&pb.RegContract_Member{
				Addr:   txutil.NewAddressFromHash(a).PBMessage(),
				Weight: int32(ratios[i]),
			})
	}

	//gen the contract addr
	data, err := newContract(addrs, ratios)
	if err != nil {
		return nil, nil, err
	}

	hash, err := hashContract(data, i.GetNonce())
	if err != nil {
		return nil, nil, err
	}

	return hash, &pb.RegContract{ContractBody: contractTx}, nil
}

func (i *GeneralCall) New_C(addrs [][]byte, ratios []int) ([]byte, error) {

	hash, msg, err := i.new(addrs, ratios)
	if err != nil {
		return nil, err
	}

	err = i.Invoke(Method_NewContract, msg)
	if err != nil {
		return nil, err
	}

	return hash, nil

}

func (i *GeneralCall) NewByDelegator(contract map[string]int32, degAddrs string) ([]byte, error) {

	var addrs [][]byte
	var ratios []int

	for a, weight := range contract {

		caddr, err := txutil.NewAddressFromString(a)
		if err != nil {
			return nil, err
		}

		addrs = append(addrs, caddr.Internal())
		ratios = append(ratios, int(weight))
	}

	hash, msg, err := i.new(addrs, ratios)
	if err != nil {
		return nil, err
	}

	if degAddrs != "" {
		degAddr, err := txutil.NewAddressFromString(degAddrs)
		if err != nil {
			return nil, err
		}
		msg.Delegator = degAddr.PBMessage()
	}

	err = i.Invoke(Method_NewContract, msg)
	if err != nil {
		return nil, err
	}

	return hash, nil

}

func (i *GeneralCall) New(contract map[string]int32) ([]byte, error) {

	return i.NewByDelegator(contract, "")
}

func (i *GeneralCall) Redeem_C(conaddr []byte, amount *big.Int, redeemAddrs [][]byte) (*pb.RedeemResponse, error) {

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

func (i *GeneralCall) Redeem(conaddr string, amount *big.Int, redeemAddrs []string) (*pb.RedeemResponse, error) {

	cconaddr, err := txutil.NewAddressFromString(conaddr)
	if err != nil {
		return nil, err
	}

	var craddrs [][]byte
	for _, addr := range redeemAddrs {
		caddr, err := txutil.NewAddressFromString(addr)
		if err != nil {
			return nil, err
		}
		craddrs = append(craddrs, caddr.Internal())
	}

	return i.Redeem_C(cconaddr.Internal(), amount, craddrs)
}

func (i *GeneralCall) Query_C(addr []byte) (error, *pb.Contract_s) {
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

func (i *GeneralCall) Query(addr string) (error, *pb.Contract_s) {

	caddr, err := txutil.NewAddressFromString(addr)
	if err != nil {
		return err, nil
	}

	return i.Query_C(caddr.Internal())
}

func (i *GeneralCall) QueryOne_C(conaddr []byte, addr []byte) (error, *pb.Contract_s) {
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

func (i *GeneralCall) QueryOne(conaddr, addr string) (error, *pb.Contract_s) {

	ccaddr, err := txutil.NewAddressFromString(conaddr)
	if err != nil {
		return err, nil
	}

	caddr, err := txutil.NewAddressFromString(addr)
	if err != nil {
		return err, nil
	}

	return i.QueryOne_C(ccaddr.Internal(), caddr.Internal())
}
