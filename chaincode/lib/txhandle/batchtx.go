package tx

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"hyperledger.abchain.org/chaincode/shim"
	txutil "hyperledger.abchain.org/core/tx"
	pb "hyperledger.abchain.org/protos"
)

//batchTx Handler integrating mutiple tx handlers and itself is another tx handlers

type batchTx struct {
	msg      pb.TxBatch
	Handlers map[string]*ChaincodeTx
}

func BatchTxHandler(hs map[string]*ChaincodeTx) *batchTx {
	return &batchTx{Handlers: hs}
}

func (h *batchTx) Msg() proto.Message { return &h.msg }

func (h *batchTx) Call(stub shim.ChaincodeStubInterface, parser txutil.Parser) ([]byte, error) {

	ret := new(pb.TxBatchResp)

	for _, subcall := range h.msg.Txs {
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
