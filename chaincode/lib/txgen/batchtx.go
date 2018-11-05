package tx

import (
	txutil "hyperledger.abchain.org/core/tx"
	pb "hyperledger.abchain.org/protos"
	"github.com/golang/protobuf/proto"
	"fmt"
)

type BatchTxCall struct {
	Msg pb.TxBatch
	TxGenerator
	query_ret chan []byte
}

func (t *BatchTxCall) Nonce() []byte {
	if t.nonce == nil {
		t.nonce = txutil.GenerateNonce()
	}
	return t.nonce
}

func (t *BatchTxCall) Invoke(method string, msg proto.Message) error {
	if len(t.Msg.Txs) > 0 && t.call_method == call_query){
		return nil, fmt.Errorf("Can not mixed with invoking and query")
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	t.Msg.Txs = append(t.Msg.Txs, &pb.TxBatchSubTx{method, payload})

	return nil
}

func (t *BatchTxCall) Query(method string, msg proto.Message) (chan QueryResp, error) {
	if len(t.Msg.Txs) > 0 && t.call_method != call_query){
		return nil, fmt.Errorf("Can not mixed with invoking and query")
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	if t.query_ret == nil{
		t.query_ret = make(chan []byte)
	}

	t.Msg.Txs = append(t.Msg.Txs, &pb.TxBatchSubTx{method, payload})

	return t.query_ret, nil
}


func (t *BatchTxCall) CommitBatch(method string) (error) {

	if len(t.Msg.Txs) == 0{
		return fmt.Errorf("Empty Batch is not allowed")
	}

	//do query
	if t.query_ret != nil{
		ret, err := t.Query(method, &t.Msg)
		if err != nil{
			return err
		}
	}
}

func (i *DeployTxCall) Deploy(method string) error {

	msg := new(pb.DeployTx)
	msg.InitParams = make(map[string][]byte)

	for method, v := range i.InitParams {
		payload, err := proto.Marshal(v)
		if err != nil {
			return err
		}

		msg.InitParams[method] = payload
	}

	err := i.txcall(method, msg)
	if err != nil {
		return err
	}

	_, err = i.postHandling(method, call_deploy)
	return err
}
