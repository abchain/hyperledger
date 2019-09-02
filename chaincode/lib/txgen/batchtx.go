package tx

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	txutil "hyperledger.abchain.org/core/tx"
	pb "hyperledger.abchain.org/protos"
)

type BatchTxCall struct {
	Msg pb.TxBatch
	*TxGenerator
	query_ret []chan QueryResp
}

func (t *BatchTxCall) Nonce() []byte {
	if t.nonce == nil {
		t.nonce = txutil.GenerateNonce()
	}
	return t.nonce
}

func (t *BatchTxCall) Invoke(method string, msg proto.Message) error {
	if t.query_ret != nil {
		return fmt.Errorf("Can not mixed with invoking and query")
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	t.Msg.Txs = append(t.Msg.Txs, &pb.TxBatchSubTx{Method: method, Payload: payload})

	return nil
}

func (t *BatchTxCall) Query(method string, msg proto.Message) (chan QueryResp, error) {

	if t.query_ret == nil {
		if len(t.Msg.Txs) > 0 {
			return nil, fmt.Errorf("Can not mixed with invoking and query")
		}
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	t.Msg.Txs = append(t.Msg.Txs, &pb.TxBatchSubTx{Method: method, Payload: payload})

	qret := make(chan QueryResp)
	t.query_ret = append(t.query_ret, qret)

	return qret, nil
}

func (t *BatchTxCall) CommitBatch(method string) error {

	if len(t.Msg.Txs) == 0 {
		return fmt.Errorf("Empty Batch is not allowed")
	}

	//do query
	if t.query_ret != nil {
		retc, err := t.TxGenerator.Query(method, &t.Msg)
		if err != nil {
			return err
		}

		t.resolveQueries(<-retc)
		return nil

	} else {
		return t.TxGenerator.Invoke(method, &t.Msg)
	}
}

func (t *BatchTxCall) resolveQueries(batchresp QueryResp) {

	if batchresp.ErrMsg != nil {
		t.FailAllQueries(batchresp.ErrMsg)
	} else {

		//this should be a batch resp
		respmsg := new(pb.TxBatchResp)
		if err := proto.Unmarshal(batchresp.SuccMsg, respmsg); err != nil {
			t.FailAllQueries(fmt.Errorf("Decode batch resp fail: %s", err))
			return
		}

		if len(respmsg.Response) != len(t.query_ret) {
			t.FailAllQueries(fmt.Errorf("Response is not matched"))
			return
		}

		for i, c := range t.query_ret {
			c <- QueryResp{respmsg.Response[i], nil}
		}
	}

}

func (t *BatchTxCall) FailAllQueries(err error) {
	for _, c := range t.query_ret {
		c <- QueryResp{nil, err}
	}
}
