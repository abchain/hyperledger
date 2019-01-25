package registrar

import (
	"errors"

	txhandle "hyperledger.abchain.org/chaincode/lib/txhandle"
	"hyperledger.abchain.org/chaincode/shim"
	"hyperledger.abchain.org/core/crypto"
	txutil "hyperledger.abchain.org/core/tx"
)

type regPreHandler struct {
	RegistrarConfig
	RegistrarTxExt
	txutil.AddrCredentials
}

func RegistrarPreHandler(cfg RegistrarConfig) *regPreHandler {
	return &regPreHandler{RegistrarConfig: cfg}
}

func (h *regPreHandler) Clone() txhandle.AddrVerifier {
	return &regPreHandler{RegistrarConfig: h.RegistrarConfig}
}

func (h *regPreHandler) Verify(addr *txutil.Address) error {
	if h.RegistrarTxExt == nil {
		return errors.New("Not inited")
	}

	pk := h.GetCredPubkey(*addr)

	if pk == nil {
		return errors.New("No credential for address FROM")
	}

	useRoot := false
	pkk := pk.GetRootFingerPrint()
	if len(pkk) == 0 {
		pkk = pk.Digest()
		useRoot = true
	}

	if err, regData := h.pubkey(pkk); err != nil {
		return err
	} else if !regData.Enabled {
		return errors.New("Registried pk is not enabled")
	} else if !useRoot {
		child, err := crypto.GetChildPublicKey(regData.Pk, pk.GetIndex())
		if err != nil {
			return err
		}

		if !child.IsEqual(pk) {
			return errors.New("Pk in credential is not matched with Registried pk")
		}
	}

	return nil

}

func (h regPreHandler) PreHandling(stub shim.ChaincodeStubInterface, _ string, tx txutil.Parser) error {

	h.AddrCredentials = tx.GetAddrCredential()
	if h.AddrCredentials == nil {
		return errors.New("Tx not include credentials")
	}

	h.RegistrarTxExt = h.NewTx(stub)
	return nil
}
