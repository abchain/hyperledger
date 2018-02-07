package service

import (
	"github.com/gocraft/web"
	"hyperledger.abchain.org/cases/ae/chaincode/cc"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
	reg "hyperledger.abchain.org/chaincode/registrar"
	"net/http"
)

const (
	regPkID    = "regKeyID"
	regPkIDDir = "/:regKeyID"
)

type registrar struct {
	*apiCore
	reg reg.GeneralCall
}

func (s *registrar) InitCaller(rw web.ResponseWriter,
	req *web.Request, next web.NextMiddlewareFunc) {

	s.reg = reg.GeneralCall{txgen.SimpleTxGen(chaincode.CC_NAME)}
	if offlineMode {
		s.reg.Dispatcher = ccCaller
	} else {
		http.Error(rw, "Not implied", http.StatusBadRequest)
		return
	}

	next(rw, req)
}

func (s *registrar) Reg(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received create registrar request")

	if s.activePrivk == nil {
		s.normalErrorF(rw, -100, "No account is specified")
		return
	}

	err := s.reg.AdminRegistrar(s.activePrivk.Public())
	if err != nil {
		s.normalError(rw, err)
	}

	s.normal(rw, &fundEntry{
		string(s.reg.Dispatcher.LastInvokeTxId()),
		encodeEntry(s.activePrivk.Public().RootFingerPrint),
	})

}

func (s *registrar) Query(rw web.ResponseWriter, req *web.Request) {

	key, err := decodeEntry(req.PathParams[regPkID])
	if err != nil {
		s.normalError(rw, err)
	}

	err, data := s.reg.Pubkey(key)
	if err != nil {
		s.normalError(rw, err)
	}

	s.normal(rw, data.RegTxid)
}
