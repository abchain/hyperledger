package subscription

import (
	"errors"
	"hyperledger.abchain.org/chaincode/generaltoken/nonce"
	tokenpb "hyperledger.abchain.org/chaincode/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/lib/caller"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	pb "hyperledger.abchain.org/chaincode/sharesubscription/protos"
	"hyperledger.abchain.org/crypto"
	txutil "hyperledger.abchain.org/tx"
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

func (i *GeneralCall) New(contract map[string]uint32, pk *crypto.PublicKey) ([]byte, error) {

	if len(contract) == 0 {
		return nil, errors.New("Empty contract")
	}

	addr, err := txutil.NewAddress(pk)
	if err != nil {
		return nil, err
	}

	contractTx := make([]*pb.RegContract_Member, 0, len(contract))

	for a, weight := range contract {
		addr, err := txutil.NewAddressFromString(a)
		if err != nil {
			return nil, err
		}

		contractTx = append(contractTx, &pb.RegContract_Member{addr.PBMessage(), weight})
	}

	msg := &pb.RegContract{addr.PBMessage(), contractTx}
	_, err = i.Invoke(Method_NewContract, msg)

	//gen the contract addr
	data, err := newContract(contract, pk)
	if err != nil {
		return nil, err
	}

	conaddr, err := hashContract(data, i.GetBuilder().GetNonce())
	if err != nil {
		return nil, err
	}

	return conaddr.Hash, err
}

func (i *GeneralCall) Redeem(conaddr []byte, addr []byte, amount *big.Int) ([]byte, error) {

	msg := &tokenpb.SimpleFund{
		amount.Bytes(),
		txutil.NewAddressFromHash(addr).PBMessage(),
		txutil.NewAddressFromHash(conaddr).PBMessage(),
	}

	_, err := i.Invoke(Method_Redeem, msg)
	if err != nil {
		return nil, err
	}

	return nonce.GeneralTokenNonceKey(i.GetBuilder().GetNonce(), conaddr, addr, amount.Bytes()), nil

}

func (i *GeneralCall) Query(addr []byte) (error, *pb.Contract) {

	msg := &pb.QueryContract{
		txutil.NewAddressFromHash(addr).PBMessage(),
		nil,
	}

	data, err := i.TxGenerator.Query(Method_Query, msg)
	if err != nil {
		return err, nil
	}

	d := &pb.Contract{}
	err = rpc.DecodeRPCResult(d, data)
	if err != nil {
		return err, nil
	}

	return nil, d
}

func (i *GeneralCall) QueryOne(conaddr []byte, addr []byte) (error, *pb.Contract) {
	msg := &pb.QueryContract{
		txutil.NewAddressFromHash(conaddr).PBMessage(),
		txutil.NewAddressFromHash(addr).PBMessage(),
	}

	data, err := i.TxGenerator.Query(Method_MemberQuery, msg)
	if err != nil {
		return err, nil
	}

	d := &pb.Contract{}
	err = rpc.DecodeRPCResult(d, data)
	if err != nil {
		return err, nil
	}

	return nil, d
}
