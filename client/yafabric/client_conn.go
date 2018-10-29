package client

import (
	"github.com/abchain/fabric/core/comm"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type ClientConn struct {
	C         *grpc.ClientConn
	BlockConn bool
}

// NewPeerClientConnection Returns a new grpc.ClientConn to the configured local PEER.
func newPeerClientConnection(block bool) (*grpc.ClientConn, error) {
	return newPeerClientConnectionWithAddress(block, viper.GetString("service.cliaddress"))
}

// NewPeerClientConnectionWithAddress Returns a new grpc.ClientConn to the configured PEER.
func newPeerClientConnectionWithAddress(block bool, peerAddress string) (*grpc.ClientConn, error) {
	if comm.TLSEnabledforService() {
		return comm.NewClientConnectionWithAddress(peerAddress, block, true, comm.InitTLSForPeer())
	}
	return comm.NewClientConnectionWithAddress(peerAddress, block, false, nil)
}

func (conn *ClientConn) Dialdefault() error {
	c, err := newPeerClientConnection(conn.BlockConn)
	if err != nil {
		return err
	}

	conn.C = c
	return nil
}

func (conn *ClientConn) Dial(server string) error {
	c, err := newPeerClientConnectionWithAddress(conn.BlockConn, server)
	if err != nil {
		return err
	}

	conn.C = c
	return nil
}
