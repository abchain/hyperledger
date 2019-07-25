package ccprotos

import (
	"fmt"
	txutil "hyperledger.abchain.org/core/tx"
	"strings"

	"hyperledger.abchain.org/chaincode/lib/runtime"
	"sort"
)

type Contract_s struct {
	contract_store
}

type ctaddr struct {
	Addr   string `asn1:"utf8"`
	Weight int32
}

type contract_store struct {
	Version   int32
	Threshold int32
	Addrs     []ctaddr
	Recursive int32 `json:"-"`
}

// func (n *account_store) Array2map() *map[string]int32 {

// 	res := make(map[string]int32)

// 	if n.Addr2WeightList != nil {
// 		for _, element := range n.Addr2WeightList {
// 			res[element.Addr] = element.Weight
// 		}
// 	}
// 	return &res
// }

// func (n *account_store) Map2array(m *map[string]int32) {

// 	var keys []string
// 	for k := range *m {
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)

// 	n.Addr2WeightList = nil

// 	for _, k := range keys {
// 		n.Addr2WeightList = append(n.Addr2WeightList, AddrByWeight_s{k, (*m)[k]})
// 	}

// }

// var notFound = AddrByWeight_s{}

// func (n *account_store) FindAndAccess(addr string) (AddrByWeight_s, int) {

// 	for i, v := range n.Addr2WeightList {
// 		if v.Addr == addr {
// 			return v, i
// 		}
// 	}
// 	return notFound, -1
// }

//notice this never check the duplication of addrs
func (n *Contract_s) Participate(addr string, weight int32) {
	n.Addrs = append(n.Addrs, ctaddr{Addr: addr, Weight: weight})
}

func (n *Contract_s) Find(addr string) int {
	for i, v := range n.Addrs {
		if addr == v.Addr {
			return i
		}
	}

	return -1
}

func (n *Contract_s) ToMap() map[string]int32 {

	r := make(map[string]int32)

	for _, v := range n.Addrs {
		r[v.Addr] = v.Weight
	}

	return r
}

func (n *Contract_s) Sorter() addrSort {
	return addrSort(n.Addrs)
}

type addrSort []ctaddr

func (s addrSort) Len() int      { return len(s) }
func (s addrSort) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s addrSort) Less(i, j int) bool {
	return strings.Compare(s[i].Addr, s[j].Addr) < 0
}

func (n *Contract_s) GetObject() interface{} { return &n.contract_store }
func (n *Contract_s) Load(interface{}) error { return nil }
func (n *Contract_s) Save() interface{}      { return n.contract_store }

func (n *Contract_s) LoadFromPB(p *Contract) *Contract_s {
	n.Threshold = p.Threshold

	for _, element := range p.Addrs {
		n.Participate(txutil.NewAddressFromHash(element.GetAddr().GetHash()).ToString(),
			element.GetWeight())
	}

	sort.Sort(n.Sorter())

	return n
}

func (n *Contract_s) ToPB() *Contract {
	res := &Contract{Threshold: n.Threshold}

	for _, element := range n.Addrs {

		addr, err := txutil.NewAddressFromString(element.Addr)
		if err != nil {
			return nil
		}
		res.Addrs = append(res.Addrs,
			&AddrByWeight{addr.PBMessage(), element.Weight})
	}
	return res
}

func (n *Contract_s) Serialize() ([]byte, error) {
	return runtime.SeralizeObject(n)
}

func (n *Contract_s) Dump() (dump string) {

	dump = dump + fmt.Sprintf("Multisign info. Threshold: [%d]\n", n.Threshold)
	for idx, e := range n.Addrs {
		dump = dump + fmt.Sprintf("		%d: [%s][%d]\n", idx, e.Addr, e.Weight)
	}

	return
}
