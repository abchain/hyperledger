package ccprotos

import (
	"github.com/golang/protobuf/proto"
	evts "hyperledger.abchain.org/chaincode/lib/events"
	txutil "hyperledger.abchain.org/core/tx"
	"math/big"
)

const (
	fundEvent   = "TRANSFERTOKEN"
	mtokenEvent = "MULTITOKENCTX"
	transResult = "TRANSFERRESULT"
)

type transferEvt func(nc []byte) (string, []byte)
type transferResult func(nc, faddr, taddr []byte,
	from *AccountData_s, to *AccountData_s) (string, []byte)
type mtokenEvt func(nc []byte, tkn string) (string, []byte)

func (transferEvt) Msg() proto.Message {
	return new(QueryTransfer)
}

func (transferEvt) Detail(msg proto.Message) interface{} {
	return NonceKey(msg.(*QueryTransfer).GetNonce())
}

func (transferResult) Msg() proto.Message {
	return new(TransferResult)
}

func newStatusOut(st *TransferResult_AccountStatus) interface{} {
	if st == nil {
		return nil
	}

	addrs, err := txutil.NewAddressFromPBMessage(st.GetAddr())
	if err != nil {
		return nil
	}

	return struct {
		Addr    string   `json:"addr"`
		Balance *big.Int `json:"balance"`
	}{
		Addr:    addrs.ToString(),
		Balance: big.NewInt(0).SetBytes(st.GetBalance()),
	}
}

func (transferResult) Detail(msg proto.Message) interface{} {

	msgR := msg.(*TransferResult)

	return struct {
		Nonce NonceKey    `json:"nonce"`
		From  interface{} `json:"from,omitempty"`
		To    interface{} `json:"to,omitempty"`
	}{msgR.GetNoncekey(), newStatusOut(msgR.GetFrom()), newStatusOut(msgR.GetTo())}
}

func (mtokenEvt) Msg() proto.Message {
	return new(MultiTokenRecord)
}

func (mtokenEvt) Detail(msg proto.Message) interface{} {
	mrec := msg.(*MultiTokenRecord)
	return struct {
		Nonce     NonceKey `json:"nonce"`
		TokenName string   `json:"token"`
	}{mrec.GetNoncekey(), mrec.GetTokenName()}
}

//NewTransferEvent is the global implement
var NewTransferEvent transferEvt

//NewTransferResult is the global implement
var NewTransferResult transferResult

//NewMultiTokenContext is the global implement
var NewMultiTokenContext mtokenEvt

func init() {

	NewTransferEvent = func(key []byte) (string, []byte) {
		bt, _ := proto.Marshal(&QueryTransfer{
			Nonce: key,
		})

		return fundEvent, bt
	}

	zeroBalance := big.NewInt(0).Bytes()

	toResultAccount := func(acc *AccountData_s, addr []byte) *TransferResult_AccountStatus {
		if acc == nil {
			return nil
		}

		var bln []byte
		if acc.Balance != nil {
			bln = acc.Balance.Bytes()
		} else {
			bln = zeroBalance
		}

		return &TransferResult_AccountStatus{
			Addr:    txutil.NewAddressFromHash(addr).PBMessage(),
			Balance: bln,
		}
	}

	NewTransferResult = func(nc, faddr, taddr []byte,
		from *AccountData_s, to *AccountData_s) (string, []byte) {

		bt, _ := proto.Marshal(&TransferResult{
			Noncekey: nc,
			From:     toResultAccount(from, faddr),
			To:       toResultAccount(to, taddr),
		})

		return transResult, bt
	}

	NewMultiTokenContext = func(nc []byte, tkn string) (string, []byte) {

		bt, _ := proto.Marshal(&MultiTokenRecord{
			Noncekey:  nc,
			TokenName: tkn,
		})

		return mtokenEvent, bt

	}

	evts.RegTxEventParser(fundEvent, NewTransferEvent)
	evts.RegTxEventParser(transResult, NewTransferResult)
	evts.RegTxEventParser(mtokenEvent, NewMultiTokenContext)

}
