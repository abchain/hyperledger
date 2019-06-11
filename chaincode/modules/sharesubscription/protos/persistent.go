package ccprotos

import (
	"hyperledger.abchain.org/core/utils"
	"hyperledger.abchain.org/protos"
	"math/big"
	"time"
)

type Contract_MemberStatus_s struct {
	Weight      int32
	TotalRedeem *big.Int
	MemberID    string `asn1:"utf8"`
}

func (n *Contract_MemberStatus_s) LoadFromPB(p *Contract_MemberStatus) {
	n.Weight = p.GetWeight()
	n.TotalRedeem = big.NewInt(0).SetBytes(p.GetTotalRedeem())
	n.MemberID = p.GetMemberID()
}

func (n *Contract_MemberStatus_s) ToPB() *Contract_MemberStatus {
	return &Contract_MemberStatus{
		Weight:      n.Weight,
		TotalRedeem: n.TotalRedeem.Bytes(),
		MemberID:    n.MemberID,
	}
}

type contract_Store struct {
	DelegatorPkFingerPrint []byte
	TotalRedeem            *big.Int
	Status                 []Contract_MemberStatus_s
	ContractTs             time.Time `asn1:"utc"`
	FrozenTo               time.Time `asn1:"utc"`
	IsFrozen               bool
	NextAddrHash           []byte
}

var notFound = Contract_MemberStatus_s{}

func (con *contract_Store) FindAndAccess(addr string) (Contract_MemberStatus_s, int) {

	for i, v := range con.Status {
		if v.MemberID == addr {
			return v, i
		}
	}

	return notFound, -1
}

func (con *contract_Store) Find(addr string) (Contract_MemberStatus_s, bool) {

	c, ind := con.FindAndAccess(addr)
	return c, ind != -1
}

type Contract_s struct {
	contract_Store
}

func (n *Contract_s) GetObject() interface{} { return &n.contract_Store }
func (n *Contract_s) Load(interface{}) error { return nil }
func (n *Contract_s) Save() interface{}      { return n.contract_Store }

func (n *Contract_s) LoadFromPB(p *Contract) {
	n.DelegatorPkFingerPrint = p.GetDelegatorPkFingerPrint()
	n.TotalRedeem = big.NewInt(0).SetBytes(p.GetTotalRedeem())
	n.IsFrozen = p.GetIsFrozen()
	n.NextAddrHash = p.GetNextAddr().GetHash()
	n.FrozenTo = utils.ConvertPBTimestamp(p.GetFrozenTo())
	for _, m := range p.GetStatus() {
		var item Contract_MemberStatus_s
		item.LoadFromPB(m)
		n.Status = append(n.Status, item)
	}
}

func (n *Contract_s) ToPB() *Contract {

	ret := &Contract{
		DelegatorPkFingerPrint: n.DelegatorPkFingerPrint,
		TotalRedeem:            n.TotalRedeem.Bytes(),
		IsFrozen:               n.IsFrozen,
		NextAddr:               &protos.TxAddr{Hash: n.NextAddrHash},
		FrozenTo:               utils.CreatePBTimestamp(n.FrozenTo),
	}

	for _, i := range n.Status {
		ret.Status = append(ret.Status, i.ToPB())
	}

	return ret
}
