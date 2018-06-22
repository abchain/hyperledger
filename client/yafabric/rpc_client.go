package client

import (
	"context"
	"github.com/abchain/fabric/peerex"
	"sync"
	"time"
)

const (
	defaultReconnectInterval = 1 * time.Minute
)

type connBuilder struct {
	sync.Mutex
	peerex.ClientConn
	waitConn      *sync.Cond
	connFail      error
	resetInterval time.Duration
}

func (c *connBuilder) close() {
	c.Lock()
	defer c.Unlock()

	if c.waitConn != nil {
		c.waitConn.Wait()
	}

	if c.C != nil {
		c.C.Close()
		c.C = nil
	}
}

func (c *connBuilder) reset(ctx context.Context) {

	to := c.resetInterval
	if int64(to) == 0 {
		to = defaultReconnectInterval
	}

	select {
	case <-time.After(to):
		c.Lock()
		c.connFail = nil
		c.Unlock()
	case <-ctx.Done():
	}
}

func (c *connBuilder) obtainConn(ctx context.Context) (*peerex.ClientConn, error) {

	c.Lock()
	defer c.Unlock()

	if c.C != nil {
		return &peerex.ClientConn{c.C, true}, nil
	}

	if c.connFail != nil {
		return nil, c.connFail
	}

	//response to do submit
	if c.waitConn == nil {
		c.waitConn = sync.NewCond(c)
		go func() {
			conn := &peerex.ClientConn{nil, true}
			err := conn.Dialdefault()

			c.Lock()
			if err != nil {
				c.connFail = err
			} else {
				c.ClientConn.C = conn.C
			}
			c.Unlock()

			c.waitConn.Broadcast()

			if err != nil {
				c.reset(ctx)
			}
		}()

		c.waitConn.Wait()
		c.waitConn = nil

	} else {
		c.waitConn.Wait()
	}

	if c.C != nil {
		return &peerex.ClientConn{c.C, true}, nil
	} else {
		return nil, c.connFail
	}

}

type RpcClientConfig struct {
	chaincodeName string
	conn          connBuilder
	security      *peerex.SecurityPolicy
	connManager   *peerex.RPCManager
	TxTimeout     time.Duration
}

func NewRPCConfig(ccName string) *RpcClientConfig {

	return &RpcClientConfig{
		chaincodeName: ccName,
		connManager:   peerex.NewRpcManager(),
	}
}

func (c *RpcClientConfig) SetUser(username string) {

	if c == nil {
		return
	}

	if c.security == nil {
		c.security = &peerex.SecurityPolicy{username, nil, nil, ""}
	} else {
		c.security.User = username
	}
}

func (c *RpcClientConfig) SetAttrs(attrs []string, isAppend bool) {

	if c == nil {
		return
	}

	if c.security == nil {
		c.security = &peerex.SecurityPolicy{"", nil, nil, ""}
	}

	if isAppend {
		c.security.Attributes = append(c.security.Attributes, attrs...)
	} else {
		c.security.Attributes = attrs
	}
}

func (c *RpcClientConfig) Quit() {

	if c == nil {
		return
	}

	c.connManager.Cancel()
	c.conn.close()
}

//adapter of the rpc caller
type rPCClient struct {
	*peerex.RpcBuilder
	lastTxid string
}

//Assign each http request (run cocurrency) a client, which can be adapted to a caller
//the client is "lazy" connect: it just do connect when required (a request has come)
//and wait for connect finish
func (c *RpcClientConfig) GetCaller() (*rPCClient, error) {

	conn, err := c.conn.obtainConn(c.connManager.Context())
	if conn == nil {
		return nil, err
	}

	builder := &peerex.RpcBuilder{
		c.chaincodeName,
		"",
		c.security,
		*conn,
		c.connManager,
		c.TxTimeout,
	}

	if err := builder.VerifyConn(); err != nil {
		return nil, err
	}

	return &rPCClient{builder, ""}, nil
}

func (r *rPCClient) LastInvokeTxId() []byte {
	return []byte(r.lastTxid)
}

func (r *rPCClient) Invoke(function string, args []string) ([]byte, error) {

	r.lastTxid = ""

	r.Function = function
	txid, err := r.Fire(args)
	if err == nil {
		r.lastTxid = txid
		return []byte(r.lastTxid), nil
	}

	return nil, err
}

func (r *rPCClient) Query(function string, args []string) ([]byte, error) {

	r.Function = function

	return r.RpcBuilder.Query(args)
}
