package subscription

import (
	"errors"
	"github.com/golang/protobuf/proto"
	pb "hyperledger.abchain.org/chaincode/modules/sharesubscription/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
)

type redeemContractAddrCred struct{}

func NewRedeemContractAddrCred(msg proto.Message) redeemContractAddrCred {

	//just testing ...
	_, ok := msg.(*pb.RedeemContract)
	if !ok {
		panic("Binding to wrong txhandler")
	}

	return redeemContractAddrCred{}
}

//this module will first catch the runtime data (as prehandler) to collect redeem address (if required), then act as verifyer,
func (h redeemContractAddrCred) PreHandling(_ shim.ChaincodeStubInterface, _ string, parser txutil.Parser) error {

	m, ok := parser.GetMessage().(*pb.RedeemContract)
	if !ok {
		return errors.New("Binding to wrong txhandler")
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

			m.Redeems = append(m.Redeems, addr.PBMessage())
		}

		parser.UpdateMsg(m)

	}

	return nil
}
