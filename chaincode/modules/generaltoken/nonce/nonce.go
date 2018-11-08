package nonce

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"hyperledger.abchain.org/chaincode/lib/runtime"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
	"math/big"
)

type TokenNonceTx interface {
	Nonce(key []byte) (error, *pb.NonceData)
	Add([]byte, *big.Int, *pb.FuncRecord, *pb.FuncRecord) error
}

type NonceConfig interface {
	NewTx(shim.ChaincodeStubInterface) TokenNonceTx
}

type StandardNonceConfig struct {
	Tag      string
	Readonly bool
}

func GeneralTokenNonceKey(txnonce []byte, from []byte, to []byte, amount []byte) []byte {

	shabyte := sha256.Sum256(bytes.Join([][]byte{txnonce,
		txutil.NewAddressFromHash(from).Hash,
		txutil.NewAddressFromHash(to).Hash,
		amount}, nil))
	return shabyte[:]
}

func nonceToKey(h []byte) string {
	return base64.RawURLEncoding.EncodeToString(h)
}

type baseNonceTx struct {
	*runtime.ChaincodeRuntime
}

const (
	nonce_tag_prefix = "GenTokenNonce_"
)

func (cfg *StandardNonceConfig) NewTx(stub shim.ChaincodeStubInterface) TokenNonceTx {
	rootname := nonce_tag_prefix + cfg.Tag

	return baseNonceTx{runtime.NewRuntime(rootname, stub, cfg.Readonly)}
}

func (nc baseNonceTx) Nonce(key []byte) (error, *pb.NonceData) {

	if len(key) != sha256.Size {
		return errors.New("Invalid nonce key length"), nil
	}

	ret := new(pb.NonceData_s)
	err := nc.Storage.Get(nonceToKey(key), ret)
	if err != nil {
		return err, nil
	}

	if ret.Txid == "" {
		return nil, nil
	}

	return nil, ret.ToPB()
}

func (nc baseNonceTx) Add(key []byte, amount *big.Int, from *pb.FuncRecord, to *pb.FuncRecord) error {

	if len(key) == 0 {
		return errors.New("Invalid (empty) key")
	}

	ret := &pb.NonceData_s{}
	err := nc.Storage.Get(nonceToKey(key), ret)
	if err != nil {
		return err
	}

	if ret.Txid != "" {
		return errors.New("Nonce is duplicated")
	}

	ret.Txid = nc.Tx.GetTxID()
	ret.Amount = amount
	ret.FromLast.LoadFromPB(from)
	ret.ToLast.LoadFromPB(to)
	txt, _ := nc.Tx.GetTxTime()
	ret.NonceTime = txt
	err = nc.Storage.Set(nonceToKey(key), ret)
	return err
}
