package subscription

import (
	"errors"
	"github.com/golang/protobuf/proto"
	pb "hyperledger.abchain.org/chaincode/modules/sharesubscription/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
)

type newContractAddrCred struct {
	*pb.RegContract
}

func NewContractAddrCred(msg proto.Message) newContractAddrCred {

	m, ok := msg.(*pb.RegContract)
	if !ok {
		panic("Binding to wrong txhandler")
	}

	return newContractAddrCred{m}
}

func (h newContractAddrCred) GetAddress() *txutil.Address {

	addr, err := txutil.NewAddressFromPBMessage(h.GetDelegatorAddr())
	if err != nil {
		return nil
	}

	return addr
}

type redeemContractAddrCred struct {
	*pb.RedeemContract

	constructMode  bool
	specifiedAddrs map[string]bool
	runtimeErr     error
}

func NewRedeemContractAddrCred(msg proto.Message) *redeemContractAddrCred {

	m, ok := msg.(*pb.RedeemContract)
	if !ok {
		panic("Binding to wrong txhandler")
	}

	return &redeemContractAddrCred{RedeemContract: m}
}

//this module will first catch the runtime data (as prehandler) to collect redeem address (if required), then act as verifyer,
func (h *redeemContractAddrCred) PreHandling(_ shim.ChaincodeStubInterface, _ string, parser txutil.Parser) error {
	if len(h.GetRedeems()) == 0 {

		cred := parser.GetAddrCredential()
		if cred == nil {
			return errors.New("Could not found which addresses should be redeem to")
		}

		for _, pk := range cred.ListCredPubkeys() {
			addr, err := txutil.NewAddress(pk)
			if err != nil {
				return err
			}

			h.Redeems = append(h.Redeems, addr.PBMessage())
		}
	}

	return nil
}

func (h *redeemContractAddrCred) ListAddress() (ret []*txutil.Address) {

	for _, redeemAddr := range h.GetRedeems() {
		if addr, err := txutil.NewAddressFromPBMessage(redeemAddr); err != nil {
			ret = append(ret, addr)
		}
	}

	return
}
