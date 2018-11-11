package subscription

import (
	"errors"
	"github.com/golang/protobuf/proto"
	pb "hyperledger.abchain.org/chaincode/modules/sharesubscription/protos"
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
	//runtime data
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

//in match mode, redeem take any address found in credentials and omit the Redeem in msg
func (h *redeemContractAddrCred) Match(addr *txutil.Address) bool {

	if len(h.GetRedeems()) == 0 {
		h.constructMode = true
	} else {
		//build matching table
		h.specifiedAddrs = make(map[string]bool)
		for _, addrpb := range h.GetRedeems() {

			addr, err := txutil.NewAddressFromPBMessage(addrpb)
			if err != nil {
				h.runtimeErr = err
				return false
			}
			h.specifiedAddrs[addr.ToString()] = true
		}
	}

	if h.constructMode {
		h.Redeems = append(h.Redeems, addr.PBMessage())
		return true
	} else {
		addrs := addr.ToString()
		_, ok := h.specifiedAddrs[addrs]
		delete(h.specifiedAddrs, addrs)
		return ok
	}
}

func (h *redeemContractAddrCred) Next(last bool) bool { return h.runtimeErr == nil }

func (h *redeemContractAddrCred) Final() error {
	if h.runtimeErr != nil {
		return h.runtimeErr
	}

	if len(h.specifiedAddrs) != 0 {
		return errors.New("No enough creds")
	}

	return nil
}
