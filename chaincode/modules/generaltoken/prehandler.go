package generaltoken

import (
	"github.com/golang/protobuf/proto"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	txutil "hyperledger.abchain.org/core/tx"
)

type fundAddrCred struct {
	*pb.SimpleFund
}

func FundAddrCred(msg proto.Message) fundAddrCred {

	m, ok := msg.(*pb.SimpleFund)
	if !ok {
		panic("Binding to wrong txhandler")
	}

	return fundAddrCred{m}
}

//and set it as RegistrarPreHandler for registrar
func (m fundAddrCred) GetAddress() *txutil.Address {

	addr, err := txutil.NewAddressFromPBMessage(m.From)
	if err != nil {
		return nil
	}

	return addr
}
