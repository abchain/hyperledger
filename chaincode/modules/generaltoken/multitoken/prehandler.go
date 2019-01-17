package multitoken

import (
	"hyperledger.abchain.org/chaincode/modules/generaltoken"
	txutil "hyperledger.abchain.org/core/tx"
)

type fundAddrCred struct {
	bind *basehandler
}

func FundAddrCred(bindh *basehandler) fundAddrCred {
	return fundAddrCred{bindh}

}

//and set it as RegistrarPreHandler for registrar
func (m fundAddrCred) GetAddress() *txutil.Address {

	m.bind.innerPrehandled()

	innerC := generaltoken.FundAddrCred(m.bind.inner.Msg())

	return innerC.GetAddress()
}
