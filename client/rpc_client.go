package client

import (
	"fmt"
	"github.com/abchain/fabric/peerex"
	log "github.com/abchain/fabric/peerex/logging"
)

type RPCClient struct {
	b *peerex.RpcBuilder
	c *peerex.ClientConn
}

var logger = log.InitLogger("CLIENT")

func NewRPCClient() (*RPCClient, error) {

	c := &RPCClient{}

	return c, nil
}

func (r *RPCClient) Connect(server string) error {

	r.c = &peerex.ClientConn{}

	logger.Debugf("Connect to gRPC server: %v", server)
	err := r.c.Dial(server) //peer.tls.rootcert.file, peer.tls.serverhostoverride
	if err != nil {
		return err
	}

	rpc := &peerex.Rpc{}

	r.b = &peerex.RpcBuilder{
		Conn:        *r.c,
		ConnManager: rpc.NewManager(),
	}

	return nil
}

func (r *RPCClient) Invoke(function string, args []string) (string, error) {

	if r.b == nil {
		return "", fmt.Errorf("RpcBuilder is nil")
	}

	r.b.Function = function

	return r.b.Fire(args)
}

func (r *RPCClient) Query(function string, args []string) ([]byte, error) {

	if r.b == nil {
		return nil, fmt.Errorf("RpcBuilder is nil")
	}

	r.b.Function = function

	return r.b.Query(args)
}

func (r *RPCClient) SetSecurityPolicy(username string) error {

	if r.b == nil {
		return fmt.Errorf("RpcBuilder is nil")
	}

	// logger.Debugf("Set fabric client user: %v", username)

	// r.b.Security = &peerex.SecurityPolicy{User: username,
	// 	Attributes: []string{security.PrivilegeAttr, security.RegionAttr}}

	return fmt.Errorf("Deprecated")
}

func (r *RPCClient) SetChaincodeName(name string) error {

	if r.b == nil {
		return fmt.Errorf("RpcBuilder is nil")
	}

	logger.Debugf("Set fabric chaincode name: %v", name)

	r.b.ChaincodeName = name

	return nil
}

func (r *RPCClient) Close() {
	if r.b != nil && r.b.Conn.C != nil {
		r.b.Conn.C.Close()
	}
}
