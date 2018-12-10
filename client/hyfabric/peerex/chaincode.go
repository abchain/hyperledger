package peerex

import (
	"context"
	"fmt"

	mspex "hyperledger.abchain.org/client/hyfabric/msp"

	fmsp "github.com/hyperledger/fabric/msp"
	fcommon "github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"
	protoutils "github.com/hyperledger/fabric/protos/utils"

	"github.com/pkg/errors"
)

const (
	query  = "query"
	invoke = "invoke"

	errorStatus = 400
)

type RPCManager struct {
	ctx    context.Context
	cancel context.CancelFunc
}

//depecrated
type RPC struct{}

func (_ *RPC) NewManager() *RPCManager {
	return NewRpcManager()
}

func NewRpcManager() *RPCManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &RPCManager{ctx, cancel}
}

func (m *RPCManager) Context() context.Context {
	return m.ctx
}

func (m *RPCManager) Cancel() {
	m.cancel()
}

//Query 查询  格式: args:[]string{"a"} 代表查询a的值 跟方法名要匹配
func (r *RPCBuilder) Query(args [][]byte) ([]byte, error) {

	peer := r.Peers[0]
	signer, err := mspex.GetSigningIdentity()
	if err != nil {
		return nil, errors.WithMessage(err, "error getting default signer")
	}

	signedProp, _, _, err := r.ChaincodeEnv.creatProposal(signer, args)
	if err != nil {
		return nil, err
	}

	ctx, cancel := r.context()
	defer cancel()

	// all responses will be checked when the signed transaction is created.
	// for now, just set this so we check the first response's status
	proposalResp, err := peer.NewEndorserClient().ProcessProposal(ctx, signedProp)
	if err != nil {
		return nil, err
	}
	if proposalResp == nil {
		return nil, errors.New("error during query: received nil proposal response")
	}
	if proposalResp.Endorsement == nil {
		return nil, errors.Errorf("endorsement failure during query. response: %v", proposalResp.Response)
	}
	return proposalResp.Response.Payload, nil
}

func (r *RPCBuilder) Invoke(args [][]byte) (string, error) {
	signer, err := mspex.GetSigningIdentity()
	if err != nil {
		return "", errors.WithMessage(err, "error getting default signer")
	}
	signedProp, txid, prop, err := r.ChaincodeEnv.creatProposal(signer, args)

	if err != nil {
		return "", err
	}
	var responses []*pb.ProposalResponse
	ctx, cancel := r.context()
	defer cancel()
	for _, peer := range r.Peers {
		//使用grpc调用endorserClient.ProcessProposal，触发endorer执行proposal  调用invoke query
		proposalResp, err := peer.NewEndorserClient().ProcessProposal(ctx, signedProp)
		if err != nil {
			return "", errors.WithMessage(err, "error endorsing ")
		}
		responses = append(responses, proposalResp)
	}
	// all responses will be checked when the signed transaction is created.
	// for now, just set this so we check the first response's status
	proposalResp := responses[0]

	for i, res := range responses {
		fmt.Println("response:", i, ":", res.Payload, "..")
	}
	//得到proposalResponse，如果是查询类命令直接返回结果；
	//如果是执行交易类，需要对交易签名CreateSignedTx，然后调用BroadcastClient发送给orderer进行排序，返回response

	if proposalResp != nil {
		if proposalResp.Response.Status >= errorStatus {
			fmt.Println("in invoke 111")
			return "", errors.New("peer response err:" + proposalResp.Response.Message)
		}
		fmt.Println("in invoke 113")
		// assemble a signed transaction (it's an Envelope message) 对交易签名CreateSignedTx
		env, err := protoutils.CreateSignedTx(prop, signer, responses...)
		if err != nil {
			return "", errors.WithMessage(err, "could not assemble transaction")
		}
		fmt.Println("in invoke 119")
		// send the envelope for ordering  调用BroadcastClient发送给orderer进行排序
		// r.OrderEnv.NodeEnv.New
		//此处需要ctx？
		bc, err := r.OrderEnv.NewBroadcastClient()
		if err != nil {
			return "", errors.WithMessage(err, "error sending transaction")
		}
		// 发送给orderer
		if err = bc.Send(env); err != nil {
			return "", errors.WithMessage(err, "error sending transaction")
		}
		defer bc.Close()
	}
	logger.Debug("invoke get txid", txid)

	return txid, nil
}

func (c *ChaincodeEnv) creatProposal(signer fmsp.SigningIdentity, args [][]byte) (*pb.SignedProposal, string, *pb.Proposal, error) {
	var (
		tMap      map[string][]byte
		channelID = c.ChannelID
		spec      = c.getChaincodeSpec(args)
	)

	// Build the ChaincodeInvocationSpec message 创建chaincode执行描述结构，创建proposal
	// invocation := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}
	creator, err := signer.Serialize()
	if err != nil {
		return nil, "", nil, errors.WithMessage(err, fmt.Sprintf("error serializing identity for %s", signer.GetIdentifier()))
	}
	invocation := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}
	prop, txid, err := protoutils.CreateChaincodeProposalWithTxIDAndTransient(fcommon.HeaderType_ENDORSER_TRANSACTION, channelID, invocation, creator, "", tMap)
	// prop, txid, err := CreateChaincodeProposalWithTxIDAndTransient(channelID, spec, creator, tMap)
	logger.Debug(" ChaincodeInvokeOrQuery CreateChaincodeProposalWithTxIDAndTransient", txid)
	if err != nil {
		return nil, "", nil, errors.WithMessage(err, "error creating proposal")
	}

	//对proposal签名
	//signedProp, err := protoutils.GetSignedProposal(prop, cf.Signer)
	signedProp, err := GetSignedProposal(prop, signer)

	if err != nil {
		return nil, "", nil, errors.WithMessage(err, "error creating signed proposal ")
	}
	logger.Debug("ChaincodeInvokeOrQuery GetSignedProposal==== success")

	return signedProp, txid, prop, nil
}

func (r *RPCBuilder) context() (context.Context, context.CancelFunc) {

	if r.ConnManager != nil {
		if r.TxTimeout == 0 {
			return context.WithCancel(r.ConnManager.ctx)
		}
		return context.WithTimeout(r.ConnManager.ctx, r.TxTimeout)
	}
	return context.WithCancel(context.Background())
}

// func (r *RPCBuilder) ChaincodeQuery(args [][]byte) (*pb.ProposalResponse, error) {
// 	peer := r.Peers[0]
// 	c := r.ChaincodeEnv
// 	signer, err := mspex.GetSigningIdentity()
// 	if err != nil {
// 		return nil, errors.WithMessage(err, "error getting default signer")
// 	}

// 	signedProp, _, _, err := c.creatProposal(signer, args)

// 	ctx, cancel := r.context()
// 	defer cancel()
// 	// res, _, _, err := c.execute(cf, args)
// 	// all responses will be checked when the signed transaction is created.
// 	// for now, just set this so we check the first response's status

// 	proposalResp, err := peer.NewEndorserClient().ProcessProposal(ctx, signedProp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if proposalResp == nil {
// 		return nil, errors.New("error during query: received nil proposal response")
// 	}
// 	if proposalResp.Endorsement == nil {
// 		return nil, errors.Errorf("endorsement failure during query. response: %v", proposalResp.Response)
// 	}

// 	return proposalResp, nil
// }

// func (r *RPCBuilder) ChaincodeInvoke(args [][]byte) (string, error) {
// 	// all responses will be checked when the signed transaction is created.
// 	// for now, just set this so we check the first response's status
// 	// responses, txid, prop, err := c.execute(cf, args)
// 	c := r.ChaincodeEnv
// 	signer, err := mspex.GetSigningIdentity()
// 	if err != nil {
// 		return "", errors.WithMessage(err, "error getting default signer")
// 	}
// 	signedProp, txid, prop, err := c.creatProposal(signer, args)
// 	if err != nil {
// 		return "", err
// 	}
// 	var responses []*pb.ProposalResponse
// 	ctx, cancel := r.context()
// 	defer cancel()
// 	for _, peer := range r.Peers {

// 		//使用grpc调用endorserClient.ProcessProposal，触发endorer执行proposal  调用invoke query
// 		proposalResp, err := peer.NewEndorserClient().ProcessProposal(ctx, signedProp)
// 		if err != nil {
// 			return "", errors.WithMessage(err, "error endorsing ")
// 		}
// 		responses = append(responses, proposalResp)
// 	}
// 	// all responses will be checked when the signed transaction is created.
// 	// for now, just set this so we check the first response's status
// 	proposalResp := responses[0]
// 	//得到proposalResponse，如果是查询类命令直接返回结果；
// 	//如果是执行交易类，需要对交易签名CreateSignedTx，然后调用BroadcastClient发送给orderer进行排序，返回response

// 	if proposalResp != nil {
// 		if proposalResp.Response.Status >= errorStatus {
// 			return "", nil
// 		}
// 		// assemble a signed transaction (it's an Envelope message) 对交易签名CreateSignedTx
// 		env, err := protoutils.CreateSignedTx(prop, signer, responses...)
// 		if err != nil {
// 			return "", errors.WithMessage(err, "could not assemble transaction")
// 		}
// 		logger.Debug("ChaincodeInvokeOrQuery protoutils.CreateSignedTx 成功")

// 		// send the envelope for ordering  调用BroadcastClient发送给orderer进行排序
// 		// r.OrderEnv.NodeEnv.New
// 		//此处需要ctx？
// 		bc, err := r.OrderEnv.NewBroadcastClient()
// 		if err != nil {
// 			return "", errors.WithMessage(err, "error sending transaction")
// 		}
// 		// 发送给orderer
// 		if err = bc.Send(env); err != nil {
// 			return "", errors.WithMessage(err, "error sending transaction")
// 		}
// 		defer bc.Close()

// 	}
// 	logger.Debug("invoke get txid", txid)

// 	return txid, nil
// }

func (r *RPCBuilder) PreHande(invoke bool) (string, error) {
	err := InitCrypto(r.MspEnv)
	if err != nil {
		return "", err
	}
	// 建立grpc 连接
	err = r.InitConn(invoke)
	if err != nil {
		return "", err
	}
	//校验grpc 连接
	err = r.VerifyConn(invoke)
	if err != nil {
		return "", err
	}
	return "", err
}
