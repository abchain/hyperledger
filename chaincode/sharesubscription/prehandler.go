package subscription

import (
	"errors"
	"github.com/abchain/fabric/core/chaincode/shim"
	txutil "hyperledger.abchain.org/tx"
)

func (h *newContractHandler) PreHandling(_ shim.ChaincodeStubInterface, _ string, tx txutil.Parser) error {

	cred := tx.GetAddrCredential()
	if cred == nil {
		return errors.New("Should not get nil cred or it should be filtered before")
	}

	addr := h.GetAddress()
	if addr == nil {
		return errors.New("Should not get nil addr or it should be filtered before")
	}

	h.pk = cred.GetCredPubkey(*addr)
	return nil
}

func (h *RegContractMsg) GetAddress() *txutil.Address {

	addr, err := txutil.NewAddressFromPBMessage(h.msg.DelegatorAddr)
	if err != nil {
		return nil
	}

	return addr
}
