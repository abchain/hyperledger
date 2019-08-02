package util

import (
	"github.com/gocraft/web"
	txgen "hyperledger.abchain.org/chaincode/lib/txgen"
)

type TxBatch struct {
	*FabricRPCCore
	txCaller  txgen.TxCaller
	batchOuts []interface{}
}

func (t *TxBatch) AddBatchOut(v interface{}) {
	t.batchOuts = append(t.batchOuts, v)
}

func (t *TxBatch) AcquireCaller() txgen.TxCaller { return t.txCaller }

type TxBatchRouter struct {
	*web.Router
}

func CreateBatchRouter(root TxRouter, path string) TxBatchRouter {
	return TxBatchRouter{
		root.Subrouter(TxBatch{}, path),
	}
}

func (r TxBatchRouter) Init(methodName string) TxBatchRouter {

	initc := func(s *TxBatch, rw web.ResponseWriter,
		req *web.Request, next web.NextMiddlewareFunc) {

		batch := &txgen.BatchTxCall{TxGenerator: s.TxGenerator}
		s.txCaller = batch

		next(rw, req)

		//childroute has been normal executed
		if rw.StatusCode() == 0 {
			err := batch.CommitBatch(methodName)
			if err != nil {
				s.NormalError(rw, err)
				return
			}

			s.DefaultOutput(s.batchOuts)
		}
	}

	r.Middleware(initc)
	return r
}
