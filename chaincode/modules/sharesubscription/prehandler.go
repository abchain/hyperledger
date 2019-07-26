package subscription

import (
	"errors"
	pb "hyperledger.abchain.org/chaincode/modules/sharesubscription/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
)

type redeemContractAddrCred struct {
	ContractConfig
}

func NewRedeemContractAddrCred(cfg ContractConfig) redeemContractAddrCred {
	return redeemContractAddrCred{cfg}
}

//this module will first catch the runtime data (as prehandler) to collect redeem address (if required)
func (v redeemContractAddrCred) PreHandling(stub shim.ChaincodeStubInterface, _ string, parser txutil.Parser) error {

	m, ok := parser.GetMessage().(*pb.RedeemContract)
	if !ok {
		return errors.New("Binding to wrong txhandler")
	}

	rt := v.NewTx(stub, parser.GetNonce())
	err, ct := rt.Query_C(m.GetContract().GetHash())
	if err != nil {
		return err
	}

	if len(m.GetRedeems()) == 0 {

		cred := parser.GetAddrCredential()
		if cred == nil {
			return errors.New("Could not found which addresses should be redeem to")
		}

		for _, pk := range cred.ListCredPubkeys() {
			addr, err := txutil.NewAddress(pk)
			if err != nil {
				return err
			}

			if _, ok := ct.Find(addr.Internal()); ok {
				m.Redeems = append(m.Redeems, addr.PBMessage())
			}
		}

		parser.UpdateMsg(m)

	} else {
		//match
		for _, addr := range m.GetRedeems() {
			caddr, err := txutil.NewAddressFromPBMessage(addr)
			if err == nil {
				if _, ok := ct.Find(caddr.Internal()); !ok {
					return errors.New("Invalid redeem addr (not in contract)")
				}
			} else {
				return err
			}
		}
	}

	return nil
}
