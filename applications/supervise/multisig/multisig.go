package multisig

import (
	"fmt"
	"github.com/gocraft/web"
	"hyperledger.abchain.org/applications/util"
	"hyperledger.abchain.org/chaincode/modules/multisign"
	tx "hyperledger.abchain.org/core/tx"
	"strconv"
	"strings"
)

const (
	ContrcatAddr = "contractAddr"
)

type Multisign struct {
	*util.FabricRPCCore
	mauth multisign.GeneralCall
}

type MultisignRouter struct {
	*web.Router
}

func CreatMultisignRouter(root util.TxRouter, path string) MultisignRouter {
	return MultisignRouter{
		root.Subrouter(Multisign{}, path),
	}
}

func (r MultisignRouter) Init() MultisignRouter {

	r.Middleware((*Multisign).InitCaller)
	return r
}

func (r MultisignRouter) BuildRoutes() {

	r.Post("/", (*Multisign).Contract)
	r.Get("/:"+ContrcatAddr, (*Multisign).Query)
	// regRouter.Post("/audit", (*RegistrarREST).Audit)
}

func (s *Multisign) InitCaller(rw web.ResponseWriter,
	req *web.Request, next web.NextMiddlewareFunc) {

	s.mauth = multisign.GeneralCall{s.TxGenerator}
	next(rw, req)
}

type contractEntry struct {
	Txid    string `json:"txID"`
	Address string `json:"contract address"`
}

func (s *Multisign) Contract(rw web.ResponseWriter, req *web.Request) {

	contractStrs := req.PostForm["contract"]
	contract := make(map[string]int32)

	for _, str := range contractStrs {
		ret := strings.Split(str, ":")
		if len(ret) < 2 {
			s.NormalErrorF(rw, -100, "Wrong contract string")
			return
		}

		w, err := strconv.Atoi(ret[1])
		if err != nil {
			s.NormalError(rw, err)
			return
		}

		contract[ret[0]] = int32(w)
	}

	threshold := 100

	if thrstr := req.PostFormValue("threshold"); thrstr != "" {
		if w, err := strconv.Atoi(thrstr); err == nil {
			threshold = w
		}
	}

	addrH, err := s.mauth.Contract(int32(threshold), contract)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	txid, err := s.TxGenerator.Result().TxID()
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, &contractEntry{
		txid,
		tx.NewAddressFromHash(addrH).ToString(),
	})
}

func (s *Multisign) Query(rw web.ResponseWriter, req *web.Request) {

	addr := req.PathParams[ContrcatAddr]
	if addr == "" {
		s.NormalError(rw, fmt.Errorf("No address"))
		return
	}

	err, data := s.mauth.Query(addr)
	if err != nil {
		s.NormalError(rw, err)
		return
	}

	s.Normal(rw, data)
}
