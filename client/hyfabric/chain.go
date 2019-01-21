package hyfabric

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"
	futils "github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
	"hyperledger.abchain.org/client"
)

// These are function names from Invoke first parameter
// 查询链数据 利用内置的cc ：qscc
const (
	GetChainInfo       string = "GetChainInfo"
	GetBlockByNumber   string = "GetBlockByNumber"
	GetBlockByHash     string = "GetBlockByHash"
	GetTransactionByID string = "GetTransactionByID"
	GetBlockByTxID     string = "GetBlockByTxID"
	InnerCC            string = "qscc"
)

//GetTransactionByID
func (r *rPCBuilder) GetTransaction(tid string) (*client.ChainTransaction, error) {
	chainName := r.ChaincodeName
	defer func() { r.ChaincodeName = chainName }()
	r.Function = GetTransactionByID
	r.ChaincodeName = InnerCC
	arg := [][]byte{
		[]byte(r.ChannelID),
		[]byte(tid),
	}
	res, err := r.Query(GetTransactionByID, arg)
	if err != nil {
		return nil, err
	}

	trans := &pb.ProcessedTransaction{}
	err = proto.Unmarshal(res, trans)
	if err != nil {
		return nil, err
	}

	tx, err := envelopeToTrasaction(trans.TransactionEnvelope)

	//tx.TxID
	if err != nil {
		return nil, err
	}
	number, err := r.GetIndexByID(tx.TxID)
	if err == nil {
		tx.Height = number
	}
	return tx, nil
}

func (r *rPCBuilder) GetIndexByID(txid string) (int64, error) {
	chainName := r.ChaincodeName
	defer func() { r.ChaincodeName = chainName }()
	r.Function = GetBlockByTxID
	r.ChaincodeName = InnerCC
	args := [][]byte{
		[]byte(r.ChannelID),
		[]byte(txid),
	}
	res, err := r.Query(GetBlockByTxID, args)
	if err != nil {
		return 0, err
	}
	blk := &common.Block{}
	err = proto.Unmarshal(res, blk)
	if err != nil {
		return 0, err
	}
	return int64(blk.Header.Number), nil
}

func (r *rPCBuilder) GetBlock(h int64) (*client.ChainBlock, error) {
	chainName := r.ChaincodeName
	defer func() {
		r.ChaincodeName = chainName
	}()
	r.Function = GetBlockByNumber
	r.ChaincodeName = InnerCC
	height := strconv.FormatInt(h, 10)
	args := [][]byte{
		[]byte(r.ChannelID),
		[]byte(height),
	}
	res, err := r.Query(GetBlockByNumber, args)
	if err != nil {
		return nil, err
	}
	blk := &common.Block{}
	err = proto.Unmarshal(res, blk)
	if err != nil {
		return nil, err
	}
	outblk := new(client.ChainBlock)
	//转化-=----
	outblk.Hash = hex.EncodeToString(blk.Header.DataHash)
	// fmt.Sprintf("%x", blk.Header.DataHash)

	outblk.Height = h
	chaincodeEvents := blockToChainCodeEvents(blk)
	for _, et := range chaincodeEvents {
		outblk.TxEvents = append(outblk.TxEvents, eventConvert(et))
	}

	for _, data := range blk.Data.Data {
		// ret := &client.ChainTransaction{}
		envelope := &common.Envelope{}
		if err = proto.Unmarshal(data, envelope); err != nil {
			return nil, fmt.Errorf("error reconstructing envelope(%s)", err)
		}
		transaction, err := envelopeToTrasaction(envelope)
		if err != nil {
			continue
		}
		ret := transaction
		ret.Height = h
		outblk.Transactions = append(outblk.Transactions, ret)
	}
	// b, _ := json.Marshal(outblk)
	return outblk, nil
}

func (r *rPCBuilder) GetTxEvent(tid string) ([]*client.ChainTxEvents, error) {
	chainName := r.ChaincodeName
	defer func() {
		r.ChaincodeName = chainName
	}()
	r.Function = GetBlockByTxID
	r.ChaincodeName = InnerCC
	args := [][]byte{
		[]byte(r.ChannelID),
		[]byte(tid),
	}
	res, err := r.Query(GetBlockByTxID, args)
	if err != nil {
		return nil, err
	}
	blk := &common.Block{}
	err = proto.Unmarshal(res, blk)
	if err != nil {
		return nil, err
	}
	txEvents := []*client.ChainTxEvents{}
	chaincodeEvents := blockToChainCodeEvents(blk)

	for _, et := range chaincodeEvents {
		if et.TxId == tid {
			txEvents = append(txEvents, eventConvert(et))
		}
	}
	return txEvents, nil
}

func (r *rPCBuilder) GetChain() (*client.Chain, error) {
	chainName := r.ChaincodeName
	defer func() { r.ChaincodeName = chainName }()
	r.Function = GetChainInfo
	r.ChaincodeName = InnerCC
	res, err := r.Query(GetChainInfo, [][]byte{[]byte(r.ChannelID)})
	if err != nil {
		return nil, err
	}
	bi := &common.BlockchainInfo{}
	err = proto.Unmarshal(res, bi)
	if err != nil {
		return nil, err
	}

	info := &client.Chain{int64(bi.Height)}
	logger.Debug("CurrentBlockHash", hex.EncodeToString(bi.CurrentBlockHash))
	logger.Debug("PreviousBlockHash", hex.EncodeToString(bi.PreviousBlockHash))
	return info, nil
}

func envelopeToTrasaction(env *common.Envelope) (*client.ChainTransaction, error) {
	ta := &client.ChainTransaction{}
	//	//====Transaction====== []
	// Height                  int64 `json:",string"`
	// TxID, Chaincode, Method string
	// CreatedFlag             bool
	// TxArgs                  [][]byte `json:"-"`

	var err error
	if env == nil {
		return ta, errors.New("common.Envelope is nil")
	}
	payl := &common.Payload{}
	err = proto.Unmarshal(env.Payload, payl)
	if err != nil {
		return nil, err
	}
	tx := &pb.Transaction{}
	err = proto.Unmarshal(payl.Data, tx)
	if err != nil {
		return nil, err
	}

	taa := &pb.TransactionAction{}
	taa = tx.Actions[0]

	cap := &pb.ChaincodeActionPayload{}
	err = proto.Unmarshal(taa.Payload, cap)
	if err != nil {
		return nil, err
	}

	pPayl := &pb.ChaincodeProposalPayload{}
	proto.Unmarshal(cap.ChaincodeProposalPayload, pPayl)

	prop := &pb.Proposal{}
	pay, err := proto.Marshal(pPayl)
	if err != nil {
		return nil, err
	}

	prop.Payload = pay
	h, err := proto.Marshal(payl.Header)
	if err != nil {
		return nil, err
	}
	prop.Header = h

	invocation := &pb.ChaincodeInvocationSpec{}
	err = proto.Unmarshal(pPayl.Input, invocation)
	if err != nil {
		return nil, err
	}

	spec := invocation.ChaincodeSpec

	// hdr := &common.Header{}
	// hdr = payl.Header
	channelHeader := &common.ChannelHeader{}
	proto.Unmarshal(payl.Header.ChannelHeader, channelHeader)

	ta.TxID = channelHeader.TxId
	ta.TxArgs = spec.GetInput().GetArgs()[1:]
	ta.Chaincode = spec.GetChaincodeId().GetName()
	ta.Method = string(spec.GetInput().GetArgs()[0])
	ta.CreatedFlag = channelHeader.GetType() == int32(common.HeaderType_ENDORSER_TRANSACTION)
	return ta, nil
}

// blockToChainCodeEvents parses block events for chaincode events associated with individual transactions
func blockToChainCodeEvents(block *common.Block) []*pb.ChaincodeEvent {
	if block == nil || block.Data == nil || block.Data.Data == nil || len(block.Data.Data) == 0 {
		return nil
	}
	events := make([]*pb.ChaincodeEvent, 0)
	//此处应该遍历block.Data.Data？
	for _, data := range block.Data.Data {
		event, err := getChainCodeEventsByByte(data)
		if err != nil {
			continue
		}
		events = append(events, event)
	}
	return events

}

func getChainCodeEventsByByte(data []byte) (*pb.ChaincodeEvent, error) {
	// env := &common.Envelope{}
	// if err := proto.Unmarshal(data, env); err != nil {
	// 	return nil, fmt.Errorf("error reconstructing envelope(%s)", err)
	// }

	env, err := futils.GetEnvelopeFromBlock(data)
	if err != nil {
		return nil, fmt.Errorf("error reconstructing envelope(%s)", err)
	}
	// get the payload from the envelope
	payload, err := futils.GetPayload(env)
	if err != nil {
		return nil, fmt.Errorf("Could not extract payload from envelope, err %s", err)
	}

	chdr, err := futils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, fmt.Errorf("Could not extract channel header from envelope, err %s", err)
	}

	if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION {

		tx, err := futils.GetTransaction(payload.Data)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling transaction payload for block event: %s", err)
		}
		//此处应该遍历tx.Actions？
		chaincodeActionPayload, err := futils.GetChaincodeActionPayload(tx.Actions[0].Payload)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling transaction action payload for block event: %s", err)
		}
		propRespPayload, err := futils.GetProposalResponsePayload(chaincodeActionPayload.Action.ProposalResponsePayload)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling proposal response payload for block event: %s", err)
		}

		caPayload, err := futils.GetChaincodeAction(propRespPayload.Extension)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling chaincode action for block event: %s", err)
		}
		ccEvent, err := futils.GetChaincodeEvents(caPayload.Events)
		if ccEvent != nil {
			return ccEvent, nil
		}

	}
	return nil, errors.New("no HeaderType_ENDORSER_TRANSACTION type ")
}

func eventConvert(event *pb.ChaincodeEvent) *client.ChainTxEvents {
	if event == nil {
		return nil
	}
	clientEvent := &client.ChainTxEvents{}
	clientEvent.Chaincode = event.ChaincodeId
	clientEvent.Name = event.EventName
	clientEvent.Payload = event.Payload
	clientEvent.TxID = event.TxId
	return clientEvent
}
