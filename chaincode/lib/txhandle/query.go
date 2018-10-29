package tx

//A prehandler for query, check if tx is expired (NEVER USE IN INVOKING until stub can provide a timestamp with consensus!)

import (
	"errors"
	txutil "hyperledger.abchain.org/tx"
	"time"
)

type IsTxExpired bool

func (force IsTxExpired) PreHandling(_ interface{}, _ string, tx txutil.Parser) error {
	expT := tx.GetTxTime()

	if expT.IsZero() && !bool(force) {
		return nil
	}

	if time.Now().After(expT) {
		return errors.New("Tx is expired")
	}

	return nil
}
