package registrar

import (
	"errors"
	"github.com/abchain/fabric/core/chaincode/shim"
	ccpb "hyperledger.abchain.org/chaincode/registrar/protos"
	"hyperledger.abchain.org/crypto"
	txutil "hyperledger.abchain.org/tx"
)

//provide prehandlers for other tx
type CheckAddressReg interface {
	GetAddress() *txutil.Address
}

type regPreHandler struct {
	RegistrarConfig
	getter CheckAddressReg
}

func RegistrarPreHandler(cfg RegistrarConfig, getter CheckAddressReg) *regPreHandler {
	return &regPreHandler{cfg, getter}
}

func (h *regPreHandler) PreHandling(stub shim.ChaincodeStubInterface, _ string, tx txutil.Parser) error {

	cred := tx.GetAddrCredential()
	if cred == nil {
		return errors.New("Tx not include credentials")
	}

	addr := h.getter.GetAddress()
	if addr == nil {
		return errors.New("No address provided")
	}

	pk := cred.GetCredPubkey(*addr)

	if pk == nil {
		return errors.New("No credential for address FROM")
	}

	reg := h.NewTx(stub)

	err, regData := reg.Pubkey(pk.RootFingerPrint)
	if err != nil {
		return err
	}

	if !regData.Enabled {
		return errors.New("Registried pk is not enabled")
	}

	rootpk, err := crypto.PublicKeyFromPBMessage(regData.Pk)
	if err != nil {
		return err
	}

	child, err := rootpk.ChildKey(pk.Index)
	if err != nil {
		return err
	}

	if !child.IsEqual(pk) {
		return errors.New("Pk in credential is not matched with Registried pk")
	}

	return nil
}

//RegPublicKey is also a prehandler
type RegCredPreHandler struct {
	msg ccpb.RegPublicKey
}

func (h *RegCredPreHandler) PreHandling(_ shim.ChaincodeStubInterface, _ string, tx txutil.Parser) error {

	cred := tx.GetAddrCredential()

	if cred == nil {
		return errors.New("Tx contains no credentials")
	}

	pk, err := crypto.PublicKeyFromPBMessage(h.msg.Pk)
	if err != nil {
		return err
	}

	addr, err := txutil.NewAddress(pk)

	if err != nil {
		return err
	}

	return cred.Verify(*addr)
}
