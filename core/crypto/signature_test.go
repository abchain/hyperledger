package crypto

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDecodeSignature_Simple(t *testing.T) {

	sample01 := "EC:01,d0de0aaeaefad02b8bdc8a01a1b8b11c696bd3d66a2c5f10780d95b7df42645cd85228a6fb29940e858e7e55842ae2bd115d1ed7cc0e82d934e929c97648cb0a,d93dd1d90fbf15d9e7df58a497aabb9b68c22fbbe7987f04538ff8aa2a26935d3cfd05f5df25798a698c8b84b47d462021b5587e753d2bc7e67258984f2ae15b:"

	s, err := DecodeCompactSignature(sample01)

	if err != nil {
		t.Fatal(err)
	}

	ecs := s.GetEc()
	if ecs == nil {
		t.Fatal("Not ecdsa")
	}

	t.Log(ecs)

	if ecs.Curvetype != 1 {
		t.Fatal("fail curvetype")
	}

	var cmpV []byte
	fmt.Sscanf("d0de0aaeaefad02b8bdc8a01a1b8b11c696bd3d66a2c5f10780d95b7df42645c", "%x", &cmpV)
	if bytes.Compare(ecs.GetP().X, cmpV) != 0 {
		t.Fatal("fail Point x", ecs.GetP().X, cmpV)
	}

	fmt.Sscanf("d85228a6fb29940e858e7e55842ae2bd115d1ed7cc0e82d934e929c97648cb0a", "%x", &cmpV)
	if bytes.Compare(ecs.GetP().Y, cmpV) != 0 {
		t.Fatal("fail Point y", ecs.GetP().Y, cmpV)
	}

	fmt.Sscanf("d93dd1d90fbf15d9e7df58a497aabb9b68c22fbbe7987f04538ff8aa2a26935d", "%x", &cmpV)
	if bytes.Compare(ecs.R, cmpV) != 0 {
		t.Fatal("fail sig R", ecs.GetR(), cmpV)
	}

	fmt.Sscanf("3cfd05f5df25798a698c8b84b47d462021b5587e753d2bc7e67258984f2ae15b", "%x", &cmpV)
	if bytes.Compare(ecs.S, cmpV) != 0 {
		t.Fatal("fail sig S", ecs.GetS(), cmpV)
	}
}
