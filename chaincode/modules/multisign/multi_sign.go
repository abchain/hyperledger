package multisign

import (
	"bytes"
	"encoding/base64"
	"errors"

	pb "hyperledger.abchain.org/chaincode/modules/multisign/protos"
	tx "hyperledger.abchain.org/core/tx"
	"hyperledger.abchain.org/core/utils"
	"sort"
)

func addrToKey(h []byte) string {
	return base64.RawURLEncoding.EncodeToString(h)
}

func hashContract(ct *pb.Contract_s, nonce []byte) ([]byte, error) {

	bt, err := ct.Serialize()
	if err != nil {
		return nil, err
	}

	bt = bytes.Join([][]byte{bt, []byte("MULTISIGN_CONTRACT")}, nil)

	ctHash, err := utils.HMACSHA256(bt, nonce)
	if err != nil {
		return nil, err
	} else if len(ctHash) > tx.ADDRESS_HASH_LEN {
		ctHash = ctHash[:tx.ADDRESS_HASH_LEN]
	}

	return ctHash, nil
}

var defaultInvokingAddrLimit = 12

func (cn *baseMultiSignTx) Contract_C(threshold int32, addrs [][]byte, weights []int32) ([]byte, error) {

	if threshold == 0 {
		return nil, errors.New("Invalid threshold")
	} else if len(addrs) != len(weights) {
		return nil, errors.New("Wrong arguments")
	}

	var invloved int
	//verifying involved contract address
	for _, kaddr := range addrs {

		_, data := cn.Query_C(kaddr)
		if data != nil {
			invloved += int(data.Recursive + 1)
			if invloved > defaultInvokingAddrLimit {
				return nil, errors.New("Involved too many contract address")
			}
		}
	}

	ct := new(pb.Contract_s)
	ct.Threshold = threshold

	for i, v := range addrs {
		ct.Participate(v, weights[i])
	}

	//sort them
	sort.Sort(ct.Sorter())

	ctHash, err := hashContract(ct, cn.nonce)
	conKey := addrToKey(ctHash)

	if retbt, _ := cn.Storage.GetRaw(conKey); len(retbt) > 0 {
		return nil, errors.New("Contract has existed")
	}

	//caution: we set recursive information here, so it is not involved in hashing,
	//which made the calling-adapter more easy to calculated a identifical one
	ct.Recursive = int32(invloved)
	//and we re-marshal the object and save it
	err = cn.Storage.Set(conKey, ct)
	if err != nil {
		return nil, err
	}

	return ctHash, nil
}

func (cn *baseMultiSignTx) Query_C(addr []byte) (error, *pb.Contract_s) {

	con := &pb.Contract_s{}

	err := cn.Storage.Get(addrToKey(addr), con)
	if err != nil {
		return err, nil
	} else if con.Threshold == 0 {
		return errors.New("Contract is not exist"), nil
	}

	return nil, con
}

func (cn *baseMultiSignTx) Update_C(addr, from, to []byte) error {

	err, ct := cn.Query_C(addr)
	if err != nil {
		return err
	}

	if index := ct.Find(to); index >= 0 {
		return errors.New("Update to existed addr is not allowed")
	}

	if index := ct.Find(from); index < 0 {
		return errors.New("Address not found")
	} else if len(to) != 0 {
		ct.Addrs[index].Addr = to
	} else {
		ct.Addrs = append(ct.Addrs[:index], ct.Addrs[index+1:]...)
	}

	return cn.Storage.Set(addrToKey(addr), ct)

}
