package ccprotos

import (
	"bytes"
	txutil "hyperledger.abchain.org/core/tx"
	"hyperledger.abchain.org/core/utils"
	"math/big"
	"time"
)

type Contract_MemberStatus_s struct {
	Weight      int32
	TotalRedeem *big.Int
	MemberID    []byte
}

func (n *Contract_MemberStatus_s) LoadFromPB(p *Contract_MemberStatus) {
	n.Weight = p.GetWeight()
	n.TotalRedeem = big.NewInt(0).SetBytes(p.GetTotalRedeem())
	n.MemberID = p.GetMemberID().GetHash()
}

func (n *Contract_MemberStatus_s) ToPB() *Contract_MemberStatus {
	return &Contract_MemberStatus{
		Weight:      n.Weight,
		TotalRedeem: n.TotalRedeem.Bytes(),
		MemberID:    txutil.NewAddressFromHash(n.MemberID).PBMessage(),
	}
}

type contract_Store struct {
	TotalWeight int64
	TotalRedeem *big.Int
	Status      []Contract_MemberStatus_s
	ContractTs  time.Time `asn1:"generalized"`
	FrozenTo    time.Time `asn1:"generalized"`
	IsFrozen    bool
}

var notFound = Contract_MemberStatus_s{}

func (con *contract_Store) FindAndAccess(addr []byte) (Contract_MemberStatus_s, int) {

	for i, v := range con.Status {
		if bytes.Compare(v.MemberID, addr) == 0 {
			return v, i
		}
	}

	return notFound, -1
}

func (con *contract_Store) Find(addr []byte) (Contract_MemberStatus_s, bool) {

	c, ind := con.FindAndAccess(addr)
	return c, ind != -1
}

type Contract_s struct {
	contract_Store
}

func (n *Contract_s) GetObject() interface{} { return &n.contract_Store }
func (n *Contract_s) Load(interface{}) error { return nil }
func (n *Contract_s) Save() interface{}      { return n.contract_Store }
func (n *Contract_s) Sorter() memberSort {
	return memberSort(n.Status)
}

type memberSort []Contract_MemberStatus_s

func (s memberSort) Len() int      { return len(s) }
func (s memberSort) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s memberSort) Less(i, j int) bool {
	return bytes.Compare(s[i].MemberID, s[j].MemberID) < 0
}

func (n *Contract_s) LoadFromPB(p *Contract) {
	n.TotalRedeem = big.NewInt(0).SetBytes(p.GetTotalRedeem())
	n.IsFrozen = p.GetIsFrozen()
	n.FrozenTo = utils.ConvertPBTimestamp(p.GetFrozenTo())
	for _, m := range p.GetStatus() {
		var item Contract_MemberStatus_s
		item.LoadFromPB(m)
		n.TotalWeight += int64(item.Weight)
		n.Status = append(n.Status, item)
	}
}

func (n *Contract_s) ToPB() *Contract {

	ret := &Contract{
		TotalRedeem: n.TotalRedeem.Bytes(),
		IsFrozen:    n.IsFrozen,
		FrozenTo:    utils.CreatePBTimestamp(n.FrozenTo),
	}

	for _, i := range n.Status {
		ret.Status = append(ret.Status, i.ToPB())
	}

	return ret
}
