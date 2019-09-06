package ccauthprotos

import (
	"hyperledger.abchain.org/chaincode/lib/runtime"
)

type Ticket_s struct {
	ticket_store
}

type ticket_store struct {
	TicketCat int32
	Dest      string `asn1:"printable"`
}

func (n *Ticket_s) GetObject() interface{} { return &n.ticket_store }
func (n *Ticket_s) Load(interface{}) error { return nil }
func (n *Ticket_s) Save() interface{}      { return n.ticket_store }

func (n *Ticket_s) LoadFromPB(p *Ticket) *Ticket_s {

	n.Dest = p.GetDesc()
	n.TicketCat = p.GetTicketCat()

	return n
}

func (n *Ticket_s) ToPB() *Ticket {

	if n == nil {
		return &Ticket{}
	}

	res := &Ticket{
		TicketCat: n.TicketCat,
		Desc: n.Dest,
	}

	return res
}

func (n *Ticket_s) Serialize() ([]byte, error) {
	return runtime.SeralizeObject(n)
}

