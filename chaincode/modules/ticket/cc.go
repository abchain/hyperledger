package ticket

type TicketTx interface {
	//add an ticket belong to the owner, which MUST has an unique id,
	//or error will be returned if the id is duplicated
	//ticket catalogy MUST start from 1 (0 or less is not allowed)
	Add(owner []byte, id []byte, ticketCat int, desc string) error
	//check if the ticket with specified id existed (error is returned
	//if not) and remove it
	Apply(owner []byte, id []byte) error
	//query the desc and catalogy, a catalogy less than 0 indicate the
	//ticket is not exist
	Query(owner []byte, id []byte) (int, string)
}
