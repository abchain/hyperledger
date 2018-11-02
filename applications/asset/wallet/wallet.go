package wallet

import (
	"bytes"
	"encoding/gob"
	"errors"
	log "github.com/op/go-logging"
	"github.com/spf13/viper"
	"hyperledger.abchain.org/core/config"
	abcrypto "hyperledger.abchain.org/core/crypto"
	"io"
	"io/ioutil"
	"sync"
)

var logger = log.MustGetLogger("WALLET")

const (
	defaultWalletFileName = "simplewallet.dat"
)

type simpleWallet struct {
	PersistFile string
	keyData     map[string]*abcrypto.PrivateKey
	lock        sync.RWMutex
}

type persistElem struct {
	Key  string
	Dump string
}

func NewWallet(fpath string) *simpleWallet {
	return &simpleWallet{
		PersistFile: fpath,
		keyData:     map[string]*abcrypto.PrivateKey{}}
}

//read path setting from viper with var "filePath"
func LoadWallet(vp *viper.Viper) *simpleWallet {

	return NewWallet(config.CanonicalizePath(vp.GetString("filePath")))
}

func (w *simpleWallet) NewPrivKey(accountID string) (*abcrypto.PrivateKey, error) {

	priv, err := abcrypto.NewPrivatekey(abcrypto.DefaultCurveType)
	if err != nil {
		return nil, err
	}

	w.lock.Lock()
	defer w.lock.Unlock()

	if _, exist := w.keyData[accountID]; exist {
		return nil, errors.New("account id exist")
	}

	w.keyData[accountID] = priv

	return priv, nil
}

func (w *simpleWallet) ImportPrivateKey(accountID string, priv *abcrypto.PrivateKey) error {

	w.lock.Lock()
	defer w.lock.Unlock()

	if _, exist := w.keyData[accountID]; exist {
		return errors.New("account id exist")
	}

	w.keyData[accountID] = priv

	return nil
}

func (w *simpleWallet) ImportPrivKey(accountID string, privkey string) error {

	priv, err := abcrypto.PrivatekeyFromString(privkey)
	if err != nil {
		return err
	}

	w.lock.Lock()
	defer w.lock.Unlock()

	if _, exist := w.keyData[accountID]; exist {
		return errors.New("account id exist")
	}

	w.keyData[accountID] = priv

	return nil
}

func (w *simpleWallet) LoadPrivKey(accountID string) (*abcrypto.PrivateKey, error) {

	w.lock.RLock()
	defer w.lock.RUnlock()

	priv, exist := w.keyData[accountID]
	if !exist {
		return nil, errors.New("account id not exist")
	}

	return priv, nil
}

func (w *simpleWallet) RemovePrivKey(accountID string) error {

	w.lock.Lock()
	defer w.lock.Unlock()

	_, exist := w.keyData[accountID]
	if !exist {
		return errors.New("account id not exist")
	}

	delete(w.keyData, accountID)

	return nil
}

func (w *simpleWallet) Rename(old string, new string) error {

	priv, err := w.LoadPrivKey(old)
	if err != nil {
		return err
	}

	_, err = w.LoadPrivKey(new)
	if err == nil {
		return errors.New("new account id already exist")
	}

	err = w.ImportPrivKey(new, priv.Str())
	if err != nil {
		return err
	}

	return w.RemovePrivKey(old)
}

func (w *simpleWallet) ListAll() (map[string]*abcrypto.PrivateKey, error) {

	w.lock.RLock()
	defer w.lock.RUnlock()

	// we do a deep copy
	copiedmap := map[string]*abcrypto.PrivateKey{}

	for k, v := range w.keyData {
		copiedmap[k] = v
	}

	return copiedmap, nil
}

func (m *simpleWallet) Load() error {

	var err error

	if m.keyData == nil {
		m.keyData = map[string]*abcrypto.PrivateKey{}
	}

	origSize := len(m.keyData)

	var data []byte
	if len(m.PersistFile) == 0 {
		data, err = ioutil.ReadFile(defaultWalletFileName)
	} else {
		data, err = ioutil.ReadFile(m.PersistFile)
	}

	if err != nil {
		return nil
	}

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	v := &persistElem{}
	err = dec.Decode(v)
	for ; err == nil; err = dec.Decode(v) {
		priv, errx := abcrypto.PrivatekeyFromString(v.Dump)
		if errx != nil {
			logger.Warning("Restore privkey <", v.Key, "> fail:", errx)
			continue
		}

		m.keyData[v.Key] = priv
	}

	if err == io.EOF {
		err = nil
	}

	logger.Debugf("Restore %d keys", len(m.keyData)-origSize)
	return nil

}

func (m *simpleWallet) Persist() error {

	m.lock.RLock()
	defer m.lock.RUnlock()

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	enc := gob.NewEncoder(buf)

	var saveSize = 0
	for k, v := range m.keyData {

		err := enc.Encode(&persistElem{k, v.Str()})
		if err != nil {
			logger.Warning("Encode privkey fail", err)
			continue
		}
		saveSize++
	}

	logger.Debugf("Save %d keys", saveSize)

	if len(m.PersistFile) == 0 {
		return ioutil.WriteFile(defaultWalletFileName, buf.Bytes(), 0666)
	} else {
		return ioutil.WriteFile(m.PersistFile, buf.Bytes(), 0666)
	}
}