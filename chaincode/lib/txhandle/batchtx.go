package tx

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	rt "hyperledger.abchain.org/chaincode/lib/runtime"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
	pb "hyperledger.abchain.org/protos"
)

//batchTx Handler integrating mutiple tx handlers and itself is another tx handlers

type batchTx struct {
	Handlers map[string]*ChaincodeTx
}

func BatchTxHandler(hs map[string]*ChaincodeTx) TxHandler {
	return batchTx{hs}
}

func (h batchTx) Msg() proto.Message { return new(pb.TxBatch) }

func (h batchTx) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	ret := new(pb.TxBatchResp)
	msg := parser.GetMessage().(*pb.TxBatch)

	//adapt the stub to caching one, we require this for batchTx
	stub = rt.NewCachingStub(stub)

	for _, subcall := range msg.Txs {
		txh, ok := h.Handlers[subcall.Method]
		if !ok {
			return nil, fmt.Errorf("No correspoinding method for calling %s", subcall.Method)
		}

		//not care about init's return
		if txret, err := txh.txSubCall(stub, subcall.Method, subcall.Payload, parser); err != nil {
			return nil, fmt.Errorf("handling method [%s] fail: %s", subcall.Method, err)
		} else {
			ret.Response = append(ret.Response, txret)
		}
	}

	retbyte, err := proto.Marshal(ret)
	if err != nil {
		return nil, fmt.Errorf("Encode response fail: %s", err)
	}

	return retbyte, nil

}

type batchArgParser struct {
	ccname         string
	wrappedParsers map[string]txutil.TxArgParser
}

type batchDetail struct {
	Method string
	Detail interface{}
}

func BatchArgParser(ccn string, p map[string]txutil.TxArgParser) txutil.TxArgParser {
	return batchArgParser{ccn, p}
}

func (g batchArgParser) Msg() proto.Message { return new(pb.TxBatch) }

func (g batchArgParser) Detail(m proto.Message) interface{} {

	var out []*batchDetail
	msg := m.(*pb.TxBatch)

	for _, subcall := range msg.Txs {
		if subparser, ok := g.wrappedParsers[subcall.GetMethod()+"@"+g.ccname]; ok {

			submsg := subparser.Msg()
			if err := proto.Unmarshal(subcall.GetPayload(), submsg); err != nil {
				out = append(out, &batchDetail{subcall.GetMethod(), err})
			} else {
				out = append(out, &batchDetail{subcall.GetMethod(), subparser.Detail(submsg)})
			}

		} else {
			out = append(out, &batchDetail{subcall.GetMethod(), subcall.GetPayload()})
		}
	}

	return out

}
