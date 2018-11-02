package subscription

import (
	"errors"
	"hyperledger.abchain.org/chaincode/shim"
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

//redeem take any address found in credentials and omit the Redeem in msg
func (m *RedeemMsg) GetAddress() *txutil.Address {
	addr, err := txutil.NewAddressFromPBMessage(m.msg.Redeem)
	if err != nil {
		return nil
	}

	return addr
}

//in match mode, redeem take any address found in credentials and omit the Redeem in msg
func (h *RedeemMsg) Match(addr *txutil.Address) bool {
	h.redeemAddr = addr
	return true
}
