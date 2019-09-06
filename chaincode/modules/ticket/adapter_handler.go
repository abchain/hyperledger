package ticket

import (
	"github.com/golang/protobuf/proto"
	rpc "hyperledger.abchain.org/chaincode/lib/caller"
	pb "hyperledger.abchain.org/chaincode/modules/ticket/protos"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
)


type addHandler struct{ TicketConfig }
type applyHandler struct{ TicketConfig }
type queryHandler struct{ TicketConfig }

func AddHandler(cfg TicketConfig) addHandler {
	return addHandler{TicketConfig: cfg}
}
func ApplyHandler(cfg TicketConfig) applyHandler {
	return applyHandler{TicketConfig: cfg}
}
func QueryHandler(cfg TicketConfig) queryHandler {
	return queryHandler{TicketConfig: cfg}
}

func (h addHandler) Msg() proto.Message     { return new(pb.Ticket) }
func (h applyHandler) Msg() proto.Message   { return new(pb.ApplyTicket) }
func (h queryHandler) Msg() proto.Message   { return new(pb.QueryTicket) }


func (h addHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	msg := parser.GetMessage().(*pb.Ticket)

	addr := msg.GetOwner().GetHash()

	err := h.NewTx(stub, parser.GetNonce()).Add(addr, msg.Id, int(msg.TicketCat), msg.Desc)
	if err != nil {
		return nil, err
	}
	return []byte("done"), nil
}

func (h applyHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {
	msg := parser.GetMessage().(*pb.ApplyTicket)

	err := h.NewTx(stub, parser.GetNonce()).Apply(msg.GetOwner().GetHash(), msg.Id)
	if err != nil {
		return nil, err
	}
	return []byte("done"), nil
}

func (h queryHandler) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error)  {

	msg := parser.GetMessage().(*pb.QueryTicket)
	catalogy, desc := h.NewTx(stub, parser.GetNonce()).Query(msg.GetOwner().GetHash(), msg.Id)

	data := &pb.Ticket_s{}
	data.TicketCat = int32(catalogy)
	data.Dest = desc
	return rpc.EncodeRPCResult(data.ToPB())
}

