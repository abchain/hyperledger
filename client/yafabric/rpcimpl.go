package client

import (
	"context"
	"fmt"
	pb "github.com/abchain/fabric/protos"
	"google.golang.org/grpc"
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

func makeStringArgsToPb(funcname string, args []string) *pb.ChaincodeInput {

	input := &pb.ChaincodeInput{}
	//please remember the trick fabric used:
	//it push the "function name" as the first argument
	//in a rpc call
	var inarg [][]byte
	if len(funcname) == 0 {
		input.Args = make([][]byte, len(args))
		inarg = input.Args[:]
	} else {
		input.Args = make([][]byte, len(args)+1)
		input.Args[0] = []byte(funcname)
		inarg = input.Args[1:]
	}

	for i, arg := range args {
		inarg[i] = []byte(arg)
	}

	return input
}

func (b *RpcBuilder) VerifyConn() error {
	if b.Conn.C == nil {
		return fmt.Errorf("Conn not inited")
	}

	s := b.Conn.C.GetState()

	if s != grpc.Ready {
		return fmt.Errorf("Conn is not ready: <%s>", s)
	}

	return nil
}

func (b *RpcBuilder) prepare(args []string) *pb.ChaincodeSpec {
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

func (b *RpcBuilder) prepareInvoke(args []string) *pb.ChaincodeInvocationSpec {

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

func (b *RpcBuilder) Deploy(args []string) (*pb.ChaincodeDeploymentSpec, error) {

	ctx, cancel := b.context()
	defer cancel()
	spec, err := pb.NewDevopsClient(b.Conn.C).Deploy(ctx, b.prepare(args))

	if err != nil {
		return nil, err
	}

	return spec, nil
}

func (b *RpcBuilder) Fire(args []string) (string, error) {

	ctx, cancel := b.context()
	defer cancel()
	resp, err := pb.NewDevopsClient(b.Conn.C).Invoke(ctx, b.prepareInvoke(args))

	if err != nil {
		return "", err
	}

	return string(resp.Msg), nil
}

func (b *RpcBuilder) Query(args []string) ([]byte, error) {

	ctx, cancel := b.context()
	defer cancel()
	resp, err := pb.NewDevopsClient(b.Conn.C).Query(ctx, b.prepareInvoke(args))

	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}
