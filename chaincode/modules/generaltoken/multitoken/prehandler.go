package generaltoken

import (
	"github.com/golang/protobuf/proto"
	pb "hyperledger.abchain.org/chaincode/modules/generaltoken/protos"
	txutil "hyperledger.abchain.org/core/tx"
)

type fundAddrCred struct {
	*pb.SimpleFund
	*pb.QueryToken
}

func FundAddrCred(msg proto.Message) fundAddrCred {

	switch m := msg.(type) {
	case *pb.SimpleFund:
		return fundAddrCred{m, nil}
	case *pb.QueryToken:
		return fundAddrCred{nil, m}
	default:
		panic("Binding to wrong txhandler")
	}

}

//and set it as RegistrarPreHandler for registrar
func (m fundAddrCred) GetAddress() *txutil.Address {

	addrpb := m.GetFrom()
	if m.SimpleFund == nil {
		addrpb = m.GetAddr()
	}

	addr, err := txutil.NewAddressFromPBMessage(addrpb)
	if err != nil {
		return nil
	}

	return addr
}
