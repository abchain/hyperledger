package registrar

import (
	"errors"

	txhandle "hyperledger.abchain.org/chaincode/lib/txhandle"
	"hyperledger.abchain.org/chaincode/shim"
	"hyperledger.abchain.org/core/crypto"
	txutil "hyperledger.abchain.org/core/tx"
)

type regPreHandler struct {
	txhandle.ListAddresses
	RegistrarConfig
}

func RegistrarPreHandler(cfg RegistrarConfig, getter txhandle.ListAddresses) *regPreHandler {
	return &regPreHandler{getter, cfg}
}

func (h *regPreHandler) PreHandling(stub shim.ChaincodeStubInterface, _ string, tx txutil.Parser) error {

	cred := tx.GetAddrCredential()
	if cred == nil {
		return errors.New("Tx not include credentials")
	}

	var addrs []*txutil.Address
	if h.ListAddresses != nil {
		addrs = h.ListAddresses(tx.GetMessage())
	} else if addrm, ok := tx.GetMessage().(txhandle.MsgAddresses); ok {
		addrs = addrm.GetAddresses()
	}

	if len(addrs) == 0 {
		return errors.New("No address provided")
	}

	reg := h.NewTx(stub)
	for _, addr := range addrs {
		pk := cred.GetCredPubkey(*addr)

		if pk == nil {
			return errors.New("No credential for address FROM")
		}

		useRoot := false
		pkk := pk.GetRootFingerPrint()
		if len(pkk) == 0 {
			pkk = pk.Digest()
			useRoot = true
		}

		if err, regData := reg.pubkey(pkk); err != nil {
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

	}

	return nil
}
