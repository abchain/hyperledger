package registar

import (
	"fmt"
	"github.com/gocraft/web"
	"hyperledger.abchain.org/applications/util"
	reg "hyperledger.abchain.org/chaincode/modules/registrar"
	"hyperledger.abchain.org/core/crypto"
)

const (
	RegPkID = "regKeyID"
)

type Registrar struct {
	*util.FabricRPCCore
	reg reg.GeneralCall
}

type RegistrarRouter struct {
	*web.Router
}

func CreatRegistrarRouter(root util.TxRouter, path string) RegistrarRouter {
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
	r.Post("/init", (*Registrar).InitReg)
	r.Get("/:"+RegPkID, (*Registrar).Query)
	// regRouter.Post("/audit", (*RegistrarREST).Audit)
}

func (s *Registrar) InitCaller(rw web.ResponseWriter,
	req *web.Request, next web.NextMiddlewareFunc) {

	s.reg = reg.GeneralCall{s.TxGenerator}
	next(rw, req)
}

type ReqEntry struct {
	Txid  string `json:"txID"`
	Nonce []byte `json:"Nonce"`
}

func (s *Registrar) InitReg(rw web.ResponseWriter, req *web.Request) {

	manager := req.FormValue("Admin")
	if manager == "" {
		manager = "Admin"
	}
	regmanager := req.FormValue("RegManager")
	if regmanager == "" {
		regmanager = manager
	}

	err := s.reg.Init(true, manager, regmanager)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	txid, err := s.TxGenerator.Result().TxID()
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, &ReqEntry{
		txid,
		s.TxGenerator.GetBuilder().GetNonce(),
	})
}

func (s *Registrar) Reg(rw web.ResponseWriter, req *web.Request) {

	logger.Debug("Received create registrar request")

	var pkbytes []byte
	var err error
	if s.ActivePrivk == nil {
		_, err = fmt.Sscanf(req.PostFormValue("publicKey"), "%x", &pkbytes)

	} else {
		pkbytes, err = crypto.PublicKeyToBytes(s.ActivePrivk.Public())
	}
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	err = s.reg.AdminRegistrar(pkbytes)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	txid, err := s.TxGenerator.Result().TxID()
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, &ReqEntry{
		string(txid),
		s.TxGenerator.GetBuilder().GetNonce(),
	})

}

func (s *Registrar) Query(rw web.ResponseWriter, req *web.Request) {

	key, err := s.DecodeEntry(req.PathParams[RegPkID])
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	err, data := s.reg.Pubkey(key)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, data.RegTxid)
}
