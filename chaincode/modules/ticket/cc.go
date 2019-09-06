package ticket

import (
	"hyperledger.abchain.org/chaincode/lib/runtime"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	"hyperledger.abchain.org/chaincode/shim"
)

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

type TicketConfig interface {
	NewTx(shim.ChaincodeStubInterface, []byte) TicketTx
}

type StandardTicketConfig struct {
	Root string
	*runtime.Config
}

func NewConfig(tag string) *StandardTicketConfig {
	cfg := runtime.NewConfig()

	return &StandardTicketConfig{ticket_tag_prefix + tag, cfg}
}

type baseTicketTx struct {
	*runtime.ChaincodeRuntime
	nonce    []byte
}


const (
	ticket_tag_prefix = "Ticket_"
)

func (cfg *StandardTicketConfig) NewTx(stub shim.ChaincodeStubInterface, nonce []byte) TicketTx {
	return &baseTicketTx{runtime.NewRuntime(cfg.Root, stub, cfg.Config), nonce}
}

type InnerInvokeConfig struct {
	txgen.InnerChaincode
}

func (c InnerInvokeConfig) NewTx(stub shim.ChaincodeStubInterface, nc []byte) TicketTx {
	return &GeneralCall{c.NewInnerTxInterface(stub, nc)}
}

