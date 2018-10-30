package nonce

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	pb "hyperledger.abchain.org/chaincode/generaltoken/protos"
	"hyperledger.abchain.org/chaincode/lib/state"
	txutil "hyperledger.abchain.org/tx"
)

type TokenNonceTx interface {
	Nonce(key []byte) (error, *pb.NonceData)
	Add([]byte, *pb.NonceData) error
}

type NonceConfig interface {
	NewTx(interface{}) TokenNonceTx
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
	state.StateMap
}

const (
	nonce_tag_prefix = "GenTokenNonce_"
)

func (cfg *StandardNonceConfig) NewTx(stub interface{}) TokenNonceTx {
	rootname := nonce_tag_prefix + cfg.Tag

	return baseNonceTx{state.NewShimMap(rootname, stub, cfg.Readonly)}
}

func (nc baseNonceTx) Nonce(key []byte) (error, *pb.NonceData) {

	if len(key) != sha256.Size {
		return errors.New("Invalid nonce key length"), nil
	}

	ret := &pb.NonceData{}

	err := nc.Get(nonceToKey(key), ret)
	if err != nil {
		return err, nil
	}

	if ret.Txid == "" {
		return nil, nil
	}

	return nil, ret
}

func (nc baseNonceTx) Add(key []byte, data *pb.NonceData) error {

	ret := &pb.NonceData{}
	err := nc.Get(nonceToKey(key), ret)
	if err != nil {
		return err
	}

	if ret.Txid != "" {
		return errors.New("Nonce is duplicated")
	}

	return nc.Set(nonceToKey(key), data)
}