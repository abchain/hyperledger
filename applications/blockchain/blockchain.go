package blockchain

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	pbempty "github.com/golang/protobuf/ptypes/empty"
	log "github.com/op/go-logging"
	"hyperledger.abchain.org/client"
	"hyperledger.abchain.org/core/tx"
	"strings"
)

var logger = log.MustGetLogger("server/blockchain")

type ChainTransaction struct {
	*client.ChainTransaction
	//Data for the original protobuf input (Message part) and Detail left for parser
	ChaincodeModule, Nonce string
	Detail, Data           interface{} `json:",omitempty"`
}

type ChainTxEvents struct {
	*client.ChainTxEvents
	Detail, Data interface{} `json:",omitempty"`
}

type ChainBlock struct {
	*client.ChainBlock
	Transactions []*ChainTransaction `json:",omitempty"`
	TxEvents     []*ChainTxEvents    `json:",omitempty"`
}

var registryParsers = map[string]abchainTx.TxArgParser{}

func handleTransaction(tx *client.ChainTransaction) *ChainTransaction {

	ret := &ChainTransaction{tx, "", "", nil, nil}

	if len(tx.TxArgs) < 2 {
		ret.Detail = notHyperledgerTx
		return ret
	}

	parser, err := abchainTx.ParseTx(new(pbempty.Empty), tx.Method, tx.TxArgs)
	if err != nil {
		ret.Detail = notHyperledgerTx
		return ret
	}
	ret.Nonce = fmt.Sprintf("%X", parser.GetNounce())
	ret.ChaincodeModule = parser.GetCCname()

	if addParser, ok := registryParsers[strings.Join([]string{ret.Method, ret.ChaincodeModule}, "@")]; ok {
		//a hack: the message is always in args[1]
		msg := addParser.Msg()
		err = proto.Unmarshal(tx.TxArgs[1], msg)
		if err != nil {
			ret.Detail = fmt.Sprintf("Invalid message arguments (%s)", err)
			return ret
		}
		ret.Data = msg
		ret.Detail = addParser.Detail(msg)
	} else {
		ret.Data = tx.TxArgs[1]
		ret.Detail = noParser
	}
	return ret

}

func handleTxEvent(txe *client.ChainTxEvents) *ChainTxEvents {

	ret := &ChainTxEvents{txe, nil, nil}

	if addParser, ok := registryParsers[txe.Name]; ok {
		//a hack: the message is always in args[2]
		msg := addParser.Msg()
		err := proto.Unmarshal(txe.Payload, msg)
		if err != nil {
			ret.Detail = fmt.Sprintf("Invalid event payload (%s)", err)
			return ret
		}
		ret.Detail = addParser.Detail(msg)
	} else {
		ret.Data = fmt.Sprintf("%X", txe.Payload)
		ret.Detail = noParser
	}

	return ret
}
