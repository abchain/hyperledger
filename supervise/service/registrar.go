package service

import (
	"github.com/gocraft/web"
	accsrv "hyperledger.abchain.org/asset/service"
	reg "hyperledger.abchain.org/chaincode/registrar"
)

const (
	RegPkID = "regKeyID"
)

type Registrar struct {
	*accsrv.RPCCoreWithAccount
	reg reg.GeneralCall
}

type RegistrarRouter struct {
	*web.Router
}

func CreatRegistrarRouter(root accsrv.RPCAccountRouter, path string) RegistrarRouter {
	return RegistrarRouter{
		root.Subrouter(Registrar{}, path),
	}
}

func (r RegistrarRouter) Init() RegistrarRouter {

	r.Middleware((*Registrar).InitCaller)
	return r
}

func (r RegistrarRouter) BuildRoutes() {

	r.Post("/", (*Registrar).Reg)
	r.Get("/:"+RegPkID, (*Registrar).Query)
	// regRouter.Post("/audit", (*RegistrarREST).Audit)
}

func (s *Registrar) InitCaller(rw web.ResponseWriter,
	req *web.Request, next web.NextMiddlewareFunc) {

	s.reg = reg.GeneralCall{s.TxGenerator}
	next(rw, req)
}

func (s *Registrar) Reg(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received create registrar request")

	if s.ActivePrivk == nil {
		s.NormalErrorF(rw, -100, "No account is specified")
		return
	}

	err := s.reg.AdminRegistrar(s.ActivePrivk.Public())
	if err != nil {
		s.NormalError(rw, err)
	}

	s.Normal(rw, &accsrv.FundEntry{
		string(s.reg.Dispatcher.LastInvokeTxId()),
		s.EncodeEntry(s.ActivePrivk.Public().RootFingerPrint),
	})

}

func (s *Registrar) Query(rw web.ResponseWriter, req *web.Request) {

	key, err := s.DecodeEntry(req.PathParams[RegPkID])
	if err != nil {
		s.NormalError(rw, err)
	}

	err, data := s.reg.Pubkey(key)
	if err != nil {
		s.NormalError(rw, err)
	}

	s.Normal(rw, data.RegTxid)
}
