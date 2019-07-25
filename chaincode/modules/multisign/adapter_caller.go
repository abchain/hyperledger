package multisign

import (
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	ccpb "hyperledger.abchain.org/chaincode/modules/multisign/protos"
	txutil "hyperledger.abchain.org/core/tx"
	"hyperledger.abchain.org/core/utils"
	txpb "hyperledger.abchain.org/protos"
)

type GeneralCall struct {
	txgen.TxCaller
}

const (
	Method_Update   = "MULTISIGN.UPDATE"
	Method_Query    = "MULTISIGN.QUERY"
	Method_Contract = "MULTISIGN.CONTRACT"
)

func (i *GeneralCall) Contract(threshold int32, addr2weight map[string]int32) ([]byte, error) {

	ct := new(ccpb.Contract)
	ct.Threshold = threshold

	for addrstr, w := range addr2weight {

		addr, err := txutil.NewAddressFromString(addrstr)
		if err != nil {
			return nil, err
		}

		ct.Addrs = append(ct.Addrs, &ccpb.AddrByWeight{
			Weight: w,
			Addr:   addr.PBMessage(),
		})
	}

	ctS := new(ccpb.Contract_s)
	ctBytes, err := ctS.LoadFromPB(ct).Serialize()
	if err != nil {
		return nil, err
	}

	hash, err := utils.HMACSHA256(ctBytes, i.GetNonce())
	if err != nil {
		return nil, err
	} else if len(hash) > txutil.ADDRESS_HASH_LEN {
		hash = hash[:txutil.ADDRESS_HASH_LEN]
	}

	if err := i.Invoke(Method_Contract, ct); err != nil {
		return nil, err
	}

	return hash, nil
}

func (i *GeneralCall) Update(addr, from, to string) error {

	var addr3 [3]*txpb.TxAddr
	var target []string
	if to == "" {
		target = []string{addr, from}
	} else {
		target = []string{addr, from, to}
	}

	for i, s := range target {

		addr, err := txutil.NewAddressFromString(s)
		if err != nil {
			return err
		}
		addr3[i] = addr.PBMessage()
	}

	return i.Invoke(Method_Update,
		&ccpb.Update{Addr: addr3[0], From: addr3[1], To: addr3[2]})

}

func (i *GeneralCall) Query(addr string) (error, *ccpb.Contract_s) {

	txaddr, err := txutil.NewAddressFromString(addr)
	if err != nil {
		return err, nil
	}

	data, err := i.TxCaller.Query(Method_Query, txaddr.PBMessage())
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
