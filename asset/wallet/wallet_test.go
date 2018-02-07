package wallet

import (
	"crypto/rand"
	"fmt"
	abcrypto "hyperledger.abchain.org/crypto"
	"math/big"
	"os"
	"path/filepath"
	"testing"
)

type sampleSets struct {
	sampleCnt     int
	existGroup    Wallet
	nonexistGroup Wallet
}

func (s *sampleSets) prepare() error {
	//simple manager should be robust and can be used for test samples
	s.existGroup = NewWallet("")
	s.nonexistGroup = NewWallet("")

	if s.sampleCnt == 0 {
		s.sampleCnt = 5000
	}

	randlimit := big.NewInt(0xFFFFFFFF)

	for i := 0; i < s.sampleCnt; i++ {
		k, err := abcrypto.NewPrivatekey(abcrypto.DefaultCurveType)
		if err != nil {
			return err
		}

		rbn, err := rand.Int(rand.Reader, randlimit)
		if err != nil {
			return err
		}

		s.existGroup.ImportPrivKey(fmt.Sprintf("TestSet%v%v", i, rbn.Int64()), k.Str())
		s.nonexistGroup.ImportPrivKey(fmt.Sprintf("TestSetNot%v%v", i, rbn.Int64()), k.Str())
	}

	return nil
}

func testChargeProcess(t *testing.T, target Wallet, sample *sampleSets) {
	vm, err := sample.existGroup.ListAll()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range vm {
		target.ImportPrivKey(k, v.Str())
	}
}

func testStandardProcess(t *testing.T, target Wallet, sample *sampleSets) {

	vm, err := sample.existGroup.ListAll()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range vm {
		v1, err := target.LoadPrivKey(k)
		if err != nil {
			t.Fatal(err)
		}

		if !v1.IsEqualForTest(v) {
			t.Fatal("Value not match for key", k)
		}
	}

	nvm, err := sample.nonexistGroup.ListAll()
	if err != nil {
		t.Fatal(err)
	}

	for k, _ := range nvm {
		_, err := target.LoadPrivKey(k)
		if err == nil {
			t.Fatal("Load key should not exist")
		}
	}
}

func TestSimpleWallet(t *testing.T) {

	walletFile := filepath.Join(os.TempDir(), "WalletTest.dat")
	t.Logf("Use wallet file: %v", walletFile)

	test1 := NewWallet(walletFile)

	set1 := &sampleSets{sampleCnt: 1000}
	err := set1.prepare()

	if err != nil {
		t.Fatal(err)
	}

	set2 := &sampleSets{sampleCnt: 2000}
	err = set2.prepare()

	if err != nil {
		t.Fatal(err)
	}

	set3 := &sampleSets{sampleCnt: 4000}
	err = set3.prepare()

	if err != nil {
		t.Fatal(err)
	}

	testChargeProcess(t, test1, set1)
	testChargeProcess(t, test1, set2)
	testStandardProcess(t, test1, set1)
	testChargeProcess(t, test1, set3)
	testStandardProcess(t, test1, set3)
	testStandardProcess(t, test1, set2)
	testStandardProcess(t, test1, set1)

	err = test1.Persist()

	if err != nil {
		t.Fatal(err)
	}

	test2 := NewWallet(walletFile)
	err = test2.Load()
	if err != nil {
		t.Fatal(err)
	}

	testStandardProcess(t, test2, set2)
	testStandardProcess(t, test2, set3)
	testStandardProcess(t, test2, set1)

	test3 := NewWallet(walletFile)
	testChargeProcess(t, test3, set2)
	err = test3.Load()
	if err != nil {
		t.Fatal(err)
	}

	testStandardProcess(t, test3, set1)
	testStandardProcess(t, test3, set3)
	testStandardProcess(t, test3, set2)

}
