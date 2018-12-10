package peerex

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"hyperledger.abchain.org/chaincode/lib/caller"
	"hyperledger.abchain.org/client"
	"hyperledger.abchain.org/client/hyfabric/utils"
)

type RpcClientConfig struct {
	// chaincodeName string
	// conn        connBuilder
	caller      *rPCBuilder
	connManager *RPCManager
	TxTimeout   time.Duration
}

var (
	logsymbol = "cc"
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
		connManager: NewRpcManager(),
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

//Load 利用viper 加载配置  完成之后 校验了peer 跟msp 信息
//orderer没有检测，不一定是invoke 操作  配置文件参考core.yaml
func (c *RpcClientConfig) Load(vp *viper.Viper) error {
	fmt.Println("1.2 fabric load config,path:", vp.ConfigFileUsed())
	initLog(vp)
	rpc := NewRpcBuilder()
	if s := vp.GetString("chaincode"); s != "" {
		rpc.ChaincodeName = s
	}
	if s := vp.GetString("channel"); s != "" {
		rpc.ChannelID = s
	}
	fmt.Println("get node conf ChannelID:", rpc.ChannelID, "chaincode:", rpc.ChaincodeName)
	nodes := vp.GetStringSlice("peers")
	fmt.Println("peers:", nodes)
	if nodes != nil && len(nodes) > 0 {
		for _, node := range nodes {
			// peer := new(PeerEnv)
			nodeCof := getConfig(vp, node)
			rpc.Peers = append(rpc.Peers, &PeerEnv{nodeCof})
		}
	} else {
		return errors.New("没有发现peers节点配置")
	}
	fmt.Println("get orderer conf")
	//是否需要读取orderer配置
	rpc.OrderEnv = &OrderEnv{
		NodeEnv: getConfig(vp, "orderer"),
	}
	fmt.Println("get msp conf")
	rpc.MspConfigPath = vp.GetString("msp.mspConfigPath")
	rpc.MspID = vp.GetString("msp.localMspId")
	rpc.MspType = vp.GetString("msp.localMspType")

	c.caller.RPCBuilder = rpc
	fmt.Println("MspConfigPath", rpc.MspConfigPath)

	err := InitCrypto(c.caller.MspEnv)
	if err != nil {
		return err
	}
	return nil
}

func getConfig(vp *viper.Viper, pre string) *NodeEnv {
	fmt.Println("get node conf pre", pre)
	node := new(NodeEnv)
	node.Address = vp.GetString(pre + ".address")
	node.HostnameOverride = vp.GetString(pre + ".serverhostoverride")
	node.TLS = vp.GetBool(pre + ".tls")
	node.RootCertFile = vp.GetString(pre + ".rootcert")
	node.ConnTimeout = vp.GetDuration(pre + ".conntimeout")
	fmt.Println(node, "---", node.Address, node.TLS, node.RootCertFile, node.HostnameOverride)
	return node
}

//Caller Assign each http request (run cocurrency) a client, which can be adapted to a caller
//the client is "lazy" connect: it just do connect when required (a request has come)
//and wait for connect finish
func (c *RpcClientConfig) Caller(spec *client.RpcSpec) (rpc.Caller, error) {
	fmt.Println("get 1.2 fabric caller")
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

func (c *RpcClientConfig) Chain() (client.ChainInfo, error) {
	return nil, fmt.Errorf("No implement")
}

func (c *RpcClientConfig) Quit() {
	fmt.Println("get 1.2 fabric Quit")
	if c == nil {
		return
	}
	c.connManager.Cancel()
	c.caller.CloseConn()
	// c.caller.close()
}

type rPCBuilder struct {
	*RPCBuilder
}

func (r *rPCBuilder) Deploy(function string, args [][]byte) (string, error) {
	fmt.Println("get 1.2 fabric Deploy")
	return "", nil
}

func (r *rPCBuilder) Invoke(function string, args [][]byte) (string, error) {
	fmt.Println("get 1.2 fabric Invoke funcName", function)
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
		fmt.Println(err)
	}
	fmt.Println("invoke success reault:", str)
	return str, err
}

func (r *rPCBuilder) Query(function string, args [][]byte) ([]byte, error) {
	fmt.Println("get 1.2 fabric Query funcName", function)
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
		fmt.Println(err)
	}
	fmt.Println("query success reault:", string(str))
	return str, err
}
