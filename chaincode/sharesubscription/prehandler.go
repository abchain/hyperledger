package subscription

import (
	"errors"
	"github.com/abchain/fabric/core/chaincode/shim"
	tokenpb "hyperledger.abchain.org/chaincode/generaltoken/protos"
	pb "hyperledger.abchain.org/chaincode/sharesubscription/protos"
	"hyperledger.abchain.org/crypto"
	txutil "hyperledger.abchain.org/tx"
)

type redeemPreHandler struct {
	msg tokenpb.SimpleFund
}

func (h *redeemPreHandler) PreHandling(_ shim.ChaincodeStubInterface, _ string, tx txutil.Parser) error {

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

func (h *redeemPreHandler) GetAddress() *txutil.Address {

	addr, err := txutil.NewAddressFromPBMessage(h.msg.To)
	if err != nil {
		return nil
	}

	return addr
}

type newContractPreHandler struct {
	msg pb.RegContract
	pk  *crypto.PublicKey
}

func (h *newContractPreHandler) PreHandling(_ shim.ChaincodeStubInterface, _ string, tx txutil.Parser) error {

	cred := tx.GetAddrCredential()

	if cred == nil {
		return errors.New("Tx contains no credentials")
	}

	addr := h.GetAddress()

	if addr == nil {
		return errors.New("Invalid address")
	}

	err := cred.Verify(*addr)
	if err != nil {
		return err
	}

	h.pk = cred.GetCredPubkey(*addr)
	return nil
}

func (h *newContractPreHandler) GetAddress() *txutil.Address {

	addr, err := txutil.NewAddressFromPBMessage(h.msg.DelegatorAddr)
	if err != nil {
		return nil
	}

	return addr
}
