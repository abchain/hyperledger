package peerex

import (
	"sync"
	"time"

	"hyperledger.abchain.org/client/hyfabric/msp"

	"google.golang.org/grpc"
)

type NodeEnv struct {
	Address          string
	HostnameOverride string
	ConnTimeout      time.Duration

	TLS          bool // 是否启用TLS 连接节点 默认是false
	RootCertFile string

	Connect *grpc.ClientConn
	sync.Mutex
	// // ClientConn
	// *RPCBuilder
	// Connect       *grpc.ClientConn
	// endpointConf  map[string]string
	waitConn      *sync.Cond
	connFail      error
	resetInterval time.Duration
}

//OrderEnv 节点的数据
type OrderEnv struct {
	*NodeEnv
}

type PeerEnv struct {
	*NodeEnv
	// sync.Mutex
	// // endpointConf  map[string]string
	// waitConn      *sync.Cond
	// connFail      error
	// resetInterval time.Duration
}

// //PeerEnv 节点的数据
// type PeersEnv struct {
// 	Peers []*PeerEnv
// }

type ChaincodeEnv struct {
	Function      string //方法名 格式:Function :query 如果为空,但如果args的len>1 则默认是invoke  否则是query
	ChaincodeName string //
	ChannelID     string //channel 的名称
}

//RPCBuilder rpc客户端公共数据
type RPCBuilder struct {
	*ChaincodeEnv
	*msp.MspEnv
	*OrderEnv
	// *PeersEnv
	Peers       []*PeerEnv
	ConnManager *RPCManager
	TxTimeout   time.Duration
}

func NewRpcBuilder() *RPCBuilder {
	r := &RPCBuilder{}
	node := new(NodeEnv)
	r.OrderEnv = new(OrderEnv)
	r.NodeEnv = node
	r.ChaincodeEnv = new(ChaincodeEnv)
	r.MspEnv = new(msp.MspEnv)
	r.Peers = make([]*PeerEnv, 0)
	return r
}

// func (p *PeersEnv) GetPeerAddresses() []string {
// 	if p == nil || p.Peers == nil {
// 		return nil
// 	}
// 	add := []string{}
// 	for _, a := range p.Peers {
// 		add = append(add, a.Address)
// 	}
// 	return add
// }

// func (p []*PeerEnv) Add(peer *PeerEnv) {
// 	p.Peers = append(p.Peers, peer)
// }
