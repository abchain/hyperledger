package ticket

import (
	"encoding/base64"
	"errors"
	pb "hyperledger.abchain.org/chaincode/modules/ticket/protos"

)

func addrToKey(h []byte) string {
	return base64.RawURLEncoding.EncodeToString(h)
}

func (cn *baseTicketTx) Add(owner []byte, id []byte, ticketCatalog int, desc string) error {

	key := append(owner, id...)
	conKey := addrToKey(key)

	if retbt, _ := cn.Storage.GetRaw(conKey); len(retbt) > 0 {
		return errors.New("Ticket already exists!")
	}

	if ticketCatalog < 1 {
		return errors.New("Invalid ticket catalog!")
	}

	ct := new(pb.Ticket_s)
	ct.TicketCat = int32(ticketCatalog)
	ct.Dest = desc

	return cn.Storage.Set(conKey, ct)
}

func (cn *baseTicketTx) Apply(owner []byte, id []byte) error {

	key := append(owner, id...)
	conKey := addrToKey(key)

	retbt, _ := cn.Storage.GetRaw(conKey)
	if len(retbt) == 0 {
		return errors.New("Ticket does not exist")
	}

	return cn.Storage.Delete(conKey)
}

func (cn *baseTicketTx) Query(owner []byte, id []byte) (int, string) {

	ticket := &pb.Ticket_s{}
	key := append(owner, id...)
	err := cn.Storage.Get(addrToKey(key), ticket)

	if err != nil || ticket.TicketCat < 1{
		return -1, ""
	}

	return int(ticket.TicketCat), ticket.Dest
}

