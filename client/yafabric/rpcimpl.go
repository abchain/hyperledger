package client

import (
	"context"
	"fmt"
	pb_empty "github.com/golang/protobuf/ptypes/empty"
	pb_wrappers "github.com/golang/protobuf/ptypes/wrappers"
	grpc_conn "google.golang.org/grpc/connectivity"
	pb "hyperledger.abchain.org/client/yafabric/protos"
	"time"
)

type RPCManager struct {
	ctx    context.Context
	cancel context.CancelFunc
}

//depecrated
type Rpc struct {
}

func (_ *Rpc) NewManager() *RPCManager {

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

type RpcBuilder struct {
	ChaincodeName string
	//	ChaincodeLang    string
	Function string

	Security *SecurityPolicy

	Conn        ClientConn
	ConnManager *RPCManager
	TxTimeout   time.Duration
}

type SecurityPolicy struct {
	User           string
	Attributes     []string
	Metadata       []byte
	CustomIDGenAlg string
}

func makeStringArgsToPb(funcname string, args [][]byte) *pb.ChaincodeInput {

	input := &pb.ChaincodeInput{}
	//please remember the trick fabric used:
	//it push the "function name" as the first argument
	//in a rpc call

	input.Args = append(input.Args, []byte(funcname))
	//TODO: here we change the byte args into string so it can passed the old fabric 0.6 defination
	//in chaincode, we will change the chaincode interface in YA-fabric later
	for _, arg := range args {
		input.Args = append(input.Args, arg)
	}

	return input
}

func (b *RpcBuilder) VerifyConn() error {
	if b.Conn.C == nil {
		return fmt.Errorf("Conn not inited")
	}

	s := b.Conn.C.GetState()

	if s != grpc_conn.Ready {
		return fmt.Errorf("Conn is not ready: <%s>", s)
	}

	return nil
}

func (b *RpcBuilder) prepare(args [][]byte) *pb.ChaincodeSpec {
	spec := &pb.ChaincodeSpec{
		Type:        pb.ChaincodeSpec_GOLANG, //always set it as golang
		ChaincodeID: &pb.ChaincodeID{Name: b.ChaincodeName},
		CtorMsg:     makeStringArgsToPb(b.Function, args),
	}

	if b.Security != nil {
		spec.SecureContext = b.Security.User
		spec.Attributes = b.Security.Attributes
		spec.Metadata = b.Security.Metadata
	}

	//final check attributes
	if spec.Attributes == nil {
		spec.Attributes = []string{}
	}

	return spec
}

func (b *RpcBuilder) prepareInvoke(args [][]byte) *pb.ChaincodeInvocationSpec {

	spec := b.prepare(args)
	invocation := &pb.ChaincodeInvocationSpec{ChaincodeSpec: spec}

	if b.Security != nil {
		if len(b.Security.CustomIDGenAlg) != 0 {
			invocation.IdGenerationAlg = b.Security.CustomIDGenAlg
		}
	}

	return invocation
}

func (b *RpcBuilder) context() (context.Context, context.CancelFunc) {

	if b.ConnManager != nil {
		if b.TxTimeout == 0 {
			return context.WithCancel(b.ConnManager.ctx)
		} else {
			return context.WithTimeout(b.ConnManager.ctx, b.TxTimeout)
		}
	} else {
		return context.WithCancel(context.Background())
	}
}

func (b *RpcBuilder) Deploy(args [][]byte) (*pb.ChaincodeDeploymentSpec, error) {

	ctx, cancel := b.context()
	defer cancel()
	spec, err := pb.NewDevopsClient(b.Conn.C).Deploy(ctx, b.prepare(args))

	if err != nil {
		return nil, err
	}

	return spec, nil
}

func (b *RpcBuilder) Fire(args [][]byte) (string, error) {

	ctx, cancel := b.context()
	defer cancel()
	resp, err := pb.NewDevopsClient(b.Conn.C).Invoke(ctx, b.prepareInvoke(args))

	if err != nil {
		return "", err
	}

	return string(resp.Msg), nil
}

func (b *RpcBuilder) Query(args [][]byte) ([]byte, error) {

	ctx, cancel := b.context()
	defer cancel()
	resp, err := pb.NewDevopsClient(b.Conn.C).Query(ctx, b.prepareInvoke(args))

	if err != nil {
		return nil, err
	} else if resp.Status != pb.Response_SUCCESS {
		return nil, fmt.Errorf("Failure resp: %s", string(resp.Msg))
	}

	return resp.Msg, nil
}

//also imply the chainAcquire
func (b *RpcBuilder) GetCurrentBlock() (int64, error) {

	ctx, cancel := b.context()
	defer cancel()
	resp, err := pb.NewOpenchainClient(b.Conn.C).GetBlockCount(ctx, new(pb_empty.Empty))

	if err != nil {
		return 0, err
	}

	return int64(resp.GetCount()), nil
}

func (b *RpcBuilder) GetBlock(blknum int64) (*pb.Block, error) {

	ctx, cancel := b.context()
	defer cancel()
	resp, err := pb.NewOpenchainClient(b.Conn.C).GetBlockByNumber(ctx, &pb.BlockNumber{Number: uint64(blknum)})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *RpcBuilder) GetTransaction(txid string) (*pb.Transaction, error) {

	ctx, cancel := b.context()
	defer cancel()
	resp, err := pb.NewOpenchainClient(b.Conn.C).GetTransactionByID(ctx, &pb_wrappers.StringValue{Value: txid})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *RpcBuilder) GetTxIndex(string) (int64, error) {
	return 0, fmt.Errorf("No implement")
}
