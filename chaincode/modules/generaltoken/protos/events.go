package ccprotos

import (
	"github.com/golang/protobuf/proto"
	evts "hyperledger.abchain.org/chaincode/lib/events"
)

const (
	fundEvent = "TRANSFERTOKEN"
)

type transferEvt func([]byte) (string, []byte)

var NewTransferEvent transferEvt

func (transferEvt) Msg() proto.Message {
	return new(QueryTransfer)
}

func (transferEvt) Detail(msg proto.Message) interface{} {
	return NonceKey(msg.(*QueryTransfer).GetNonce())
}

func init() {

	NewTransferEvent = func(key []byte) (string, []byte) {
		bt, _ := proto.Marshal(&QueryTransfer{
			Nonce: key,
		})

		return fundEvent, bt
	}

	evts.RegTxEventParser(fundEvent, NewTransferEvent)

}
