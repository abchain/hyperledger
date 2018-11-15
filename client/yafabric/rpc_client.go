package client

import (
	"context"
	"github.com/spf13/viper"
	"hyperledger.abchain.org/chaincode/lib/caller"
	"hyperledger.abchain.org/client"
	"sync"
	"time"
)

const (
	defaultReconnectInterval = 1 * time.Minute
)

type connBuilder struct {
	sync.Mutex
	ClientConn
	endpointConf  map[string]string
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

func (c *connBuilder) obtainConn(ctx context.Context) (*ClientConn, error) {

	c.Lock()
	defer c.Unlock()

	if c.C != nil {
		return &ClientConn{c.C, true}, nil
	}

	if c.connFail != nil {
		return nil, c.connFail
	}

	//response to do submit
	if c.waitConn == nil {
		c.waitConn = sync.NewCond(c)
		go func() {
			conn := &ClientConn{nil, true}
			err := conn.Dial(c.endpointConf)

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
		return &ClientConn{c.C, true}, nil
	} else {
		return nil, c.connFail
	}

}

type RpcClientConfig struct {
	chaincodeName string
	conn          connBuilder
	security      *SecurityPolicy
	connManager   *RPCManager
	TxTimeout     time.Duration
}

func NewRPCConfig() client.RpcClient {

	return &RpcClientConfig{
		connManager: NewRpcManager(),
	}
}

func init() {
	client.Client_Impls["yafabric"] = NewRPCConfig
}

func setUser(s *SecurityPolicy, username string) *SecurityPolicy {

	if s == nil {
		s = &SecurityPolicy{username, nil, nil, ""}
	} else {
		s.User = username
	}
	return s
}

func setAttrs(s *SecurityPolicy, attrs []string, isAppend bool) *SecurityPolicy {

	if s == nil {
		s = &SecurityPolicy{"", nil, nil, ""}
	}

	if isAppend {
		s.Attributes = append(s.Attributes, attrs...)
	} else {
		s.Attributes = attrs
	}
	return s
}

/*
	the configuration for client can include these fields:
	- chaincode
	- username
	- userattr (a list of strings)
	- endpoint (for compatible, its subfield can be put in top level now but this
			    will be deprecated later)
		- server
		- tlsenabled
		- certfile
		- hostname (override the hostname in certfile)

*/
func (c *RpcClientConfig) Load(vp *viper.Viper) error {
	if s := vp.GetString("chaincode"); s != "" {
		c.chaincodeName = s
	}

	if s := vp.GetString("username"); s != "" {
		c.security = setUser(c.security, s)
	}

	if s := vp.GetStringSlice("userattr"); s != nil {
		c.security = setAttrs(c.security, s, false)
	}

	if vp.IsSet("endpoint") {
		c.conn.endpointConf = vp.GetStringMapString("endpoint")
	} else {
		c.conn.endpointConf = map[string]string{
			"server":     vp.GetString("server"),
			"tlsenabled": vp.GetString("tlsenabled"),
			"certfile":   vp.GetString("certfile"),
		}
	}

	return nil
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
	*RpcBuilder
}

//Assign each http request (run cocurrency) a client, which can be adapted to a caller
//the client is "lazy" connect: it just do connect when required (a request has come)
//and wait for connect finish
func (c *RpcClientConfig) Caller(spec *client.RpcSpec) (rpc.Caller, error) {

	conn, err := c.conn.obtainConn(c.connManager.Context())
	if conn == nil {
		return nil, err
	}

	ccname := c.chaincodeName
	sec := c.security
	if spec != nil {
		if spec.ChaincodeName != "" {
			ccname = spec.ChaincodeName
		}

		extendAttr := true
		if spec.Options != nil {
			//if the usenamem is not reset, we suppose caller wish to add
			//more attributes instead of replace it
			if s := spec.Options.GetString("username"); s != "" {
				sec = setUser(sec, s)
				extendAttr = false
			}
		}

		if len(spec.Attributes) > 0 {
			sec = setAttrs(sec, spec.Attributes, extendAttr)
		}
	}

	builder := &RpcBuilder{
		ccname,
		"",
		sec,
		*conn,
		c.connManager,
		c.TxTimeout,
	}

	if err := builder.VerifyConn(); err != nil {
		return nil, err
	}

	return &rPCClient{builder}, nil
}

func (r *rPCClient) Deploy(function string, args [][]byte) (string, error) {
	return "", nil
}

func (r *rPCClient) Invoke(function string, args [][]byte) (string, error) {

	r.Function = function
	txid, err := r.Fire(args)
	if err == nil {
		return txid, nil
	}

	return "", err
}

func (r *rPCClient) Query(function string, args [][]byte) ([]byte, error) {

	r.Function = function

	return r.RpcBuilder.Query(args)
}
