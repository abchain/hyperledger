package events

import (
	"fmt"
	tx "hyperledger.abchain.org/core/tx"
)

var eventAggregator = map[string]tx.TxArgParser{}

type EventDuplicated string

func (s EventDuplicated) Error() string {
	return fmt.Sprintf("Event <%s> is duplicated", string(s))
}

func MergeTxEventParsers(parsers map[string]tx.TxArgParser) error {

	for n, p := range eventAggregator {

		if _, existed := parsers[n]; existed {
			return EventDuplicated(n)
		} else {
			parsers[n] = p
		}
	}

	return nil

}

func MustMergeTxEventParsers(parsers map[string]tx.TxArgParser) {
	err := MergeTxEventParsers(parsers)

	if err != nil {
		panic(err)
	}
}

func RegTxEventParser(n string, p tx.TxArgParser) {

	if _, existed := eventAggregator[n]; existed {
		panic(EventDuplicated(n))
	} else {
		eventAggregator[n] = p
	}
}
