package multisign

import (
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	ccpb "hyperledger.abchain.org/chaincode/modules/multisign/protos"
	txutil "hyperledger.abchain.org/core/tx"
	txpb "hyperledger.abchain.org/protos"
	"sort"
)

type GeneralCall struct {
	txgen.TxCaller
}

const (
	Method_Update   = "MULTISIGN.UPDATE"
	Method_Query    = "MULTISIGN.QUERY"
	Method_Contract = "MULTISIGN.CONTRACT"
)

func (i *GeneralCall) Contract_C(threshold int32, addrs [][]byte, weights []int32) ([]byte, error) {

	ct := new(ccpb.Contract)
	ct.Threshold = threshold

	for i, add := range addrs {
		ct.Addrs = append(ct.Addrs, &ccpb.AddrByWeight{
			Weight: weights[i],
			Addr:   txutil.NewAddressFromHash(add).PBMessage(),
		})
	}

	ctS := new(ccpb.Contract_s)
	sort.Sort(ctS.LoadFromPB(ct).Sorter())

	hash, err := hashContract(ctS, i.GetNonce())
	if err != nil {
		return nil, err
	}

	if err := i.Invoke(Method_Contract, ct); err != nil {
		return nil, err
	}

	return hash, nil

}

func (i *GeneralCall) Contract(threshold int32, addr2weight map[string]int32) ([]byte, error) {

	addrs := make([][]byte, 0, len(addr2weight))
	weights := make([]int32, 0, len(addr2weight))

	for addr, w := range addr2weight {

		caddr, err := txutil.NewAddressFromString(addr)
		if err != nil {
			return nil, err
		}

		weights = append(weights, w)
		addrs = append(addrs, caddr.Hash)
	}

	return i.Contract_C(threshold, addrs, weights)

}

func (i *GeneralCall) Update_C(addr, from, to []byte) error {

	var addr3 [3]*txpb.TxAddr
	var target [][]byte
	if to == nil {
		target = [][]byte{addr, from}
	} else {
		target = [][]byte{addr, from, to}
	}

	for i, h := range target {
		addr3[i] = txutil.NewAddressFromHash(h).PBMessage()
	}

	return i.Invoke(Method_Update,
		&ccpb.Update{Addr: addr3[0], From: addr3[1], To: addr3[2]})

}

func (i *GeneralCall) Update(addr, from, to string) error {

	var addr3 [3]*txutil.Address
	var target []string
	if to == "" {
		target = []string{addr, from}
	} else {
		target = []string{addr, from, to}
	}

	for i, s := range target {

		var err error
		addr3[i], err = txutil.NewAddressFromString(s)
		if err != nil {
			return err
		}
	}

	return i.Update_C(addr3[0].Hash, addr3[1].Hash, addr3[2].Internal())

}

func (i *GeneralCall) Query_C(addr []byte) (error, *ccpb.Contract_s) {

	data, err := i.TxCaller.Query(Method_Query, txutil.NewAddressFromHash(addr).PBMessage())
	if err != nil {
		return err, nil
	}

	d := &ccpb.Contract{}
	err = txgen.SyncQueryResult(d, data)
	if err != nil {
		return err, nil
	}

	ret := new(ccpb.Contract_s)
	return nil, ret.LoadFromPB(d)

}

func (i *GeneralCall) Query(addr string) (error, *ccpb.Contract_s) {

	txaddr, err := txutil.NewAddressFromString(addr)
	if err != nil {
		return err, nil
	}

	return i.Query_C(txaddr.Hash)
}
