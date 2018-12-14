package hyfabric

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"hyperledger.abchain.org/chaincode/lib/caller"
	"hyperledger.abchain.org/client"
	"hyperledger.abchain.org/client/hyfabric/peerex"
	"hyperledger.abchain.org/client/hyfabric/utils"
)

type RpcClientConfig struct {
	// chaincodeName string
	// conn        connBuilder
	caller      *rPCBuilder
	connManager *peerex.RPCManager
	TxTimeout   time.Duration
}

var (
	logsymbol = "hyfabric"
	logger    = utils.MustGetLogger(logsymbol)
)

func init() {
	if client.Client_Impls == nil {
		client.Client_Impls = make(map[string]func() client.RpcClient)
	}
	client.Client_Impls["hyfabric"] = NewRPCConfig
}

// //是否跟orderer通讯
// func (c *connBuilder) obtainConn() error {
// 	return c.RPCBuilder.InitConn(true)
// }

// func (c *connBuilder) verifyConn() error {
// 	return c.RPCBuilder.VerifyConn(true)
// }

func NewRPCConfig() client.RpcClient {
	return &RpcClientConfig{
		connManager: peerex.NewRpcManager(),
		caller:      new(rPCBuilder),
	}
}

func initLog(vp *viper.Viper) {
	l := vp.GetString("logging.level")
	if !utils.IsNullOrEmpty(l) {
		utils.SetModuleLevel("^"+logsymbol, l)
	}
	format := vp.GetString("logging.format")
	if !utils.IsNullOrEmpty(format) {
		utils.SetFormat(format)
	}
}

//Load 利用viper 加载配置   配置文件参考core.yaml
func (c *RpcClientConfig) Load(vp *viper.Viper) error {
	initLog(vp)
	logger.Debug("1.2 fabric load config,path:", vp.ConfigFileUsed())
	rpc := peerex.NewRpcBuilder()
	if s := vp.GetString("chaincode"); s != "" {
		rpc.ChaincodeName = s
	}
	if s := vp.GetString("channel"); s != "" {
		rpc.ChannelID = s
	}
	logger.Debug("get node conf ChannelID:", rpc.ChannelID, "chaincode:", rpc.ChaincodeName)
	nodes := vp.GetStringSlice("peers")
	logger.Debug("peers:", nodes)
	if nodes != nil && len(nodes) > 0 {
		for _, node := range nodes {
			// peer := new(PeerEnv)
			nodeCof := getConfig(vp, node)
			rpc.Peers = append(rpc.Peers, &peerex.PeerEnv{nodeCof})
		}
	} else {
		return errors.New("没有发现peers节点配置")
	}
	logger.Debug("get orderer conf")
	//是否需要读取orderer配置
	rpc.OrderEnv = &peerex.OrderEnv{
		NodeEnv: getConfig(vp, "orderer"),
	}
	logger.Debug("get msp conf")
	rpc.MspConfigPath = vp.GetString("msp.mspConfigPath")
	rpc.MspID = vp.GetString("msp.localMspId")
	rpc.MspType = vp.GetString("msp.localMspType")

	c.caller.RPCBuilder = rpc
	logger.Debug("MspConfigPath", rpc.MspConfigPath)

	err := peerex.InitCrypto(c.caller.MspEnv)
	if err != nil {
		return err
	}
	return nil
}

func getConfig(vp *viper.Viper, pre string) *peerex.NodeEnv {
	logger.Debug("get node conf pre", pre)
	node := new(peerex.NodeEnv)
	node.Address = vp.GetString(pre + ".address")
	node.HostnameOverride = vp.GetString(pre + ".serverhostoverride")
	node.TLS = vp.GetBool(pre + ".tls")
	node.RootCertFile = vp.GetString(pre + ".rootcert")
	node.ConnTimeout = vp.GetDuration(pre + ".conntimeout")
	logger.Debug(node, "---", node.Address, node.TLS, node.RootCertFile, node.HostnameOverride)
	return node
}

//Caller Assign each http request (run cocurrency) a client, which can be adapted to a caller
//the client is "lazy" connect: it just do connect when required (a request has come)
//and wait for connect finish
func (c *RpcClientConfig) Caller(spec *client.RpcSpec) (rpc.Caller, error) {
	logger.Debug("get 1.2 fabric caller")
	//先初始化数据，之后校验数据，再进行grpc连接
	builder := c.caller.RPCBuilder
	ccname := builder.ChaincodeName
	if spec != nil {
		if spec.ChaincodeName != "" {
			ccname = spec.ChaincodeName
		}
	}
	builder.ChaincodeName = ccname
	builder.TxTimeout = c.TxTimeout
	builder.ConnManager = c.connManager

	// err := InitCrypto(c.caller.MspEnv)
	// if err != nil {
	// 	return nil, err
	// }
	return &rPCBuilder{
		RPCBuilder: builder,
	}, nil
}

//返回链信息的接口
func (c *RpcClientConfig) Chain() (client.ChainInfo, error) {
	//先初始化数据，之后校验数据，再进行grpc连接
	builder := c.caller.RPCBuilder
	builder.TxTimeout = c.TxTimeout
	builder.ConnManager = c.connManager
	return &rPCBuilder{
		RPCBuilder: builder,
	}, nil
	// return nil, fmt.Errorf("No implement")
}

func (c *RpcClientConfig) Quit() {
	logger.Debug("get 1.2 fabric Quit")
	if c == nil {
		return
	}
	c.connManager.Cancel()
	c.caller.CloseConn()
	// c.caller.close()
}

type rPCBuilder struct {
	*peerex.RPCBuilder
}

func (r *rPCBuilder) Deploy(function string, args [][]byte) (string, error) {
	logger.Debug("get 1.2 fabric Deploy")
	return "", nil
}

func (r *rPCBuilder) Invoke(function string, args [][]byte) (string, error) {
	logger.Debug("get 1.2 fabric Invoke funcName", function, "chainName:", r.ChaincodeName)
	r.Function = function
	// 建立grpc 连接
	err := r.RPCBuilder.InitConn(true)
	if err != nil {
		return "", err
	}
	//校验grpc 连接
	err = r.RPCBuilder.VerifyConn(true)
	if err != nil {
		return "", err
	}
	str, err := r.RPCBuilder.Invoke(args)

	if err != nil {
		return "", err
	}
	logger.Debug("invoke success reault:", str)
	return str, err
}

func (r *rPCBuilder) Query(function string, args [][]byte) ([]byte, error) {
	logger.Debug("get 1.2 fabric Query funcName", function, "chainName:", r.ChaincodeName)
	r.Function = function
	//建立grpc 连接
	err := r.RPCBuilder.InitConn(false)
	if err != nil {
		return nil, err
	}
	//校验grpc 连接
	err = r.RPCBuilder.VerifyConn(false)
	if err != nil {
		return nil, err
	}
	str, err := r.RPCBuilder.Query(args)

	if err != nil {
		return nil, err
	}
	logger.Debug("query success reault:", string(str))
	return str, err
}
