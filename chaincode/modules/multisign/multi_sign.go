package multisign

import (
	"errors"
	pb "hyperledger.abchain.org/chaincode/modules/multisign/protos"
	tx "hyperledger.abchain.org/core/tx"
	"hyperledger.abchain.org/core/utils"
	"sort"
)

var defaultInvokingAddrLimit = 12

func (cn *baseMultiSignTx) Contract(threshold int32, addr2weight map[string]int32) ([]byte, error) {

	if threshold == 0 {
		return nil, errors.New("Invalid threshold")
	}

	var invloved int
	//verifying involved contract address
	for kaddr, _ := range addr2weight {

		_, data := cn.Query(kaddr)
		if data != nil {
			invloved += int(data.Recursive + 1)
			if invloved > defaultInvokingAddrLimit {
				return nil, errors.New("Involved too many contract address")
			}
		}
	}

	ct := new(pb.Contract_s)
	ct.Threshold = threshold

	for k, v := range addr2weight {
		ct.Participate(k, v)
	}

	//sort them
	sort.Sort(ct.Sorter())

	bt, err := ct.Serialize()
	if err != nil {
		return nil, err
	}

	ctHash, err := utils.HMACSHA256(bt, cn.nonce)
	if err != nil {
		return nil, err
	} else if len(ctHash) > tx.ADDRESS_HASH_LEN {
		ctHash = ctHash[:tx.ADDRESS_HASH_LEN]
	}

	conAddr := tx.NewAddressFromHash(ctHash).ToString()

	if retbt, _ := cn.Storage.GetRaw(conAddr); len(retbt) > 0 {
		return nil, errors.New("Contract has existed")
	}

	//caution: we set recursive information here, so it is not involved in hashing,
	//which made the calling-adapter more easy to calculated a identifical one
	ct.Recursive = int32(invloved)
	//and we re-marshal the object and save it
	err = cn.Storage.Set(conAddr, ct)
	if err != nil {
		return nil, err
	}

	return ctHash, nil
}

func (cn *baseMultiSignTx) Query(addr string) (error, *pb.Contract_s) {

	con := &pb.Contract_s{}

	err := cn.Storage.Get(addr, con)
	if err != nil {
		return err, nil
	} else if con.Threshold == 0 {
		return errors.New("Contract is not exist"), nil
	}

	return nil, con
}

func (cn *baseMultiSignTx) Update(addr, from, to string) error {

	err, ct := cn.Query(addr)
	if err != nil {
		return err
	}

	if index := ct.Find(to); index >= 0 {
		return errors.New("Update to existed addr is not allowed")
	}

	if index := ct.Find(from); index < 0 {
		return errors.New("Address not found")
	} else if to != "" {
		ct.Addrs[index].Addr = to
	} else {
		ct.Addrs = append(ct.Addrs[:index], ct.Addrs[index+1:]...)
	}

	return cn.Storage.Set(addr, ct)

}
