package ticket

import (
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	ccpb "hyperledger.abchain.org/chaincode/modules/ticket/protos"
	pb "hyperledger.abchain.org/chaincode/modules/ticket/protos"
	txutil "hyperledger.abchain.org/core/tx"
)

type GeneralCall struct {
	txgen.TxCaller
}

const (

	Method_Add      = "TICKET.UPDATE"
	Method_Apply    = "TICKET.QUERY"
	Method_Query    = "TICKET.CONTRACT"
)

func (i *GeneralCall) Add(owner []byte, id []byte, ticketCat int, desc string) error {

	ct := new(ccpb.Ticket)
	ct.Id = id
	ct.Owner = txutil.NewAddressFromHash(owner).PBMessage()
	ct.TicketCat = int32(ticketCat)
	ct.Desc = desc


	if err := i.Invoke(Method_Add, ct); err != nil {
		return err
	}

	return nil
}


func (i *GeneralCall) Apply(owner []byte, id []byte) error {

	ct := new(ccpb.ApplyTicket)
	ct.Id = id
	ct.Owner = txutil.NewAddressFromHash(owner).PBMessage()

	return i.Invoke(Method_Apply, ct)
}

func (i *GeneralCall) Query(owner []byte, id []byte) (int, string) {

	msg := &pb.QueryTicket{
		Owner: txutil.NewAddressFromHash(owner).PBMessage(),
		Id: id,
	}

	data, err := i.TxCaller.Query(Method_Query, msg)
	if err != nil {
		return -1, ""
	}

	d := &pb.Ticket{}
	err = txgen.SyncQueryResult(d, data)
	if err != nil {
		return -1, ""
	}

	ret := new(pb.Ticket_s)
	ret.LoadFromPB(d)
	return int(ret.TicketCat), ret.Dest
}


