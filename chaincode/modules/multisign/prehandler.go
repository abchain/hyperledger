package multisign

import (
	"errors"
	"github.com/golang/protobuf/proto"

	txhandle "hyperledger.abchain.org/chaincode/lib/txhandle"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
)

type addrVerifier struct {
	MultiSignConfig
	preSetH        txhandle.TxPreHandler
	recursiveLimit int

	rt                MultiSignAddressTx
	recursiveVerifier func(*txutil.Address) error
}

var defaultRecursiveDepth = 3

func MultiSignAddrPreHandler(cfg MultiSignConfig) *addrVerifier {
	return &addrVerifier{MultiSignConfig: cfg, recursiveLimit: defaultRecursiveDepth}
}

func (h *addrVerifier) SetRecursiveDepth(d int) {
	h.recursiveLimit = d
}

//bind preset addrCredVerifier instead of a default one, notice it MUST
//not contain addrVerifier itself, or code just panic here
func (h *addrVerifier) BindPreset(ph txhandle.TxPreHandler) {

	//sanity check
	checkH := txhandle.NewAddrCredVerifierFromTemplate(nil, ph)
	if !checkH.AddVerifier(h) {
		panic("Target handler include verifier itself!")
	}

	h.preSetH = ph
}

func (h *addrVerifier) Verify(addr *txutil.Address) error {
	if h.rt == nil {
		return errors.New("Not inited")
	}

	err, contract := h.rt.Query_C(addr.Hash)
	if err != nil {
		return err
	}

	//now each address in contract is checked, and threshold for passed address
	//is accumulated until it can pass the threshold
	var thresholdacc int32

	for _, addr := range contract.Addrs {

		if err := h.recursiveVerifier(txutil.NewAddressFromHash(addr.Addr)); err != nil {
			continue
		} else {
			thresholdacc += addr.Weight
		}

		if thresholdacc >= contract.Threshold {
			//yeah, we have passed!
			return nil
		}

	}

	return errors.New("No enough credential")

}

func (h *addrVerifier) PreHandling(stub shim.ChaincodeStubInterface,
	funcname string, tx txutil.Parser) error {

	h.rt = h.NewTx(stub, tx.GetNonce())
	h.recursiveVerifier = func(addr *txutil.Address) error {
		la := func(proto.Message) []*txutil.Address {
			return []*txutil.Address{addr}
		}

		var subHandler txhandle.TxPreHandler
		if h.preSetH == nil {
			subHandler = txhandle.NewAddrCredVerifier(la)
		} else {
			subHandler = txhandle.NewAddrCredVerifierFromTemplate(la, h.preSetH)
		}

		if h.recursiveLimit > 0 {
			hcpy := &addrVerifier{
				MultiSignConfig: h.MultiSignConfig,
				preSetH:         h.preSetH,
				recursiveLimit:  h.recursiveLimit - 1,
			}

			txhandle.AttachAddrVerifier([]txhandle.TxPreHandler{subHandler}, hcpy)
		}

		return subHandler.PreHandling(stub, funcname, tx)
	}

	return nil
}
