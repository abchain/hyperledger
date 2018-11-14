package util

import (
	"strings"
	"testing"
)

func TestWindowsExpand(t *testing.T) {
	
	pathorigin0 := "d:\\1234\\5678\\"
	pathorigin1 := "d:\\1234\\5678"
	pathorigin2 := "%appdata%\\5678\\"
	pathorigin3 := "%appdata%\\5678"
	
	ret0 := CanonicalizePath(pathorigin0)
	if strings.Compare(pathorigin0, ret0) != 0{
		t.Fatalf("Fail string 0: output is <%s>", ret0)
	}
	
	ret1 := CanonicalizePath(pathorigin1)
	if strings.Compare(pathorigin0, ret1) != 0{
		t.Fatalf("Fail string 1: output is <%s>", ret1)
	}
	
	ret2 := CanonicalizePath(pathorigin2)
	if strings.Compare(pathorigin2, ret2) == 0{
		t.Fatalf("Fail string 2: output is <%s>", ret2)
	}
	
	t.Log("ret2 output: ", ret2)
	
	ret3 := CanonicalizePath(pathorigin3)
	if strings.Compare(ret2, ret3) != 0{
		t.Fatalf("Fail string 3: output is <%s>", ret3)
	}	
}


