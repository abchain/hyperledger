package generaltoken

import (
	"errors"
	"github.com/abchain/fabric/core/chaincode/shim"
	ccpb "hyperledger.abchain.org/chaincode/generaltoken/protos"
	txutil "hyperledger.abchain.org/tx"
)

//SimpleFund also has a prehandler
type FundCredPreHandler struct {
	msg ccpb.SimpleFund
}

func (h *FundCredPreHandler) PreHandling(_ shim.ChaincodeStubInterface, _ string, tx txutil.Parser) error {

	cred := tx.GetAddrCredential()

	if cred == nil {
		return errors.New("Tx contains no credentials")
	}

	addr := h.GetAddress()

	if addr == nil {
		return errors.New("Invalid address")
	}

	return cred.Verify(*addr)
}

//and set it as RegistrarPreHandler for registrar
func (h *FundCredPreHandler) GetAddress() *txutil.Address {

	addr, err := txutil.NewAddressFromPBMessage(h.msg.From)
	if err != nil {
		return nil
	}

	return addr
}
