package registrar

import (
	"errors"

	"hyperledger.abchain.org/chaincode/lib/txhandle"
	"hyperledger.abchain.org/chaincode/shim"
	"hyperledger.abchain.org/core/crypto"
	txutil "hyperledger.abchain.org/core/tx"
)

type regPreHandler struct {
	tx.ParseAddress
	RegistrarConfig
}

func RegistrarPreHandler(cfg RegistrarConfig, getter tx.ParseAddress) *regPreHandler {
	return &regPreHandler{getter, cfg}
}

func (h *regPreHandler) PreHandling(stub shim.ChaincodeStubInterface, _ string, tx txutil.Parser) error {

	cred := tx.GetAddrCredential()
	if cred == nil {
		return errors.New("Tx not include credentials")
	}

	addr := h.GetAddress()
	if addr == nil {
		return errors.New("No address provided")
	}

	pk := cred.GetCredPubkey(*addr)

	if pk == nil {
		return errors.New("No credential for address FROM")
	}

	reg := h.NewTx(stub)

	useRoot := false
	pkk := pk.GetRootFingerPrint()
	if len(pkk) == 0 {
		pkk = pk.Digest()
		useRoot = true
	}
	err, regData := reg.pubkey(pkk)
	if err != nil {
		return err
	}

	if !regData.Enabled {
		return errors.New("Registried pk is not enabled")
	}

	if useRoot {
		return nil
	}

	child, err := crypto.GetChildPublicKey(regData.Pk, pk.GetIndex())
	if err != nil {
		return err
	}

	if !child.IsEqual(pk) {
		return errors.New("Pk in credential is not matched with Registried pk")
	}

	return nil
}

func (m *RegPkMsg) GetAddress() *txutil.Address {

	pk, err := crypto.PublicKeyFromBytes(m.msg.PkBytes)
	if err != nil {
		return nil
	}

	addr, err := txutil.NewAddress(pk)

	if err != nil {
		return nil
	}

	return addr
}
