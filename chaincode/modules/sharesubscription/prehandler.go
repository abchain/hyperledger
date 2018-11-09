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

	addr, err := txutil.NewAddressFromPBMessage(h.DelegatorAddr)
	if err != nil {
		return nil
	}

	return addr
}

type redeemContractAddrCred struct {
	*pb.RedeemContract
}

func NewRedeemContractAddrCred(msg proto.Message) redeemContractAddrCred {

	m, ok := msg.(*pb.RedeemContract)
	if !ok {
		panic("Binding to wrong txhandler")
	}

	return redeemContractAddrCred{m}
}

//redeem take any address found in credentials and omit the Redeem in msg
func (m *RedeemMsg) GetAddress() *txutil.Address {
	addr, err := txutil.NewAddressFromPBMessage(m.msg.Redeem)
	if err != nil {
		return nil
	}

	return addr
}

//in match mode, redeem take any address found in credentials and omit the Redeem in msg
func (h redeemContractAddrCred) Match(addr *txutil.Address) bool {
	h.redeemAddr = addr
	return true
}

func (h redeemContractAddrCred) Next(last bool) bool {

}
