package ccauthprotos

import (
	"sort"
	"testing"
)

func TestSorting(t *testing.T) {
	ct := new(Contract_s)

	ct.Participate("bbb1", 10)
	ct.Participate("aa", 9)
	ct.Participate("zzz", 8)
	ct.Participate("ccc", 7)

	sort.Sort(ct.Sorter())
	t.Log(ct.Addrs)

	if ct.Addrs[0].Addr != "aa" || ct.Addrs[0].Weight != 9 {
		t.Fatal("Wrong sorting")
	}

	if ct.Addrs[1].Addr != "bbb1" || ct.Addrs[1].Weight != 10 {
		t.Fatal("Wrong sorting")
	}

	if ct.Addrs[2].Addr != "ccc" || ct.Addrs[2].Weight != 7 {
		t.Fatal("Wrong sorting")
	}

	if ct.Addrs[3].Addr != "zzz" || ct.Addrs[3].Weight != 8 {
		t.Fatal("Wrong sorting")
	}
}
