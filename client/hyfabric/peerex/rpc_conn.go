package peerex

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	grpc_conn "google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
)

const (
	defaultReconnectInterval = 1 * time.Minute
)

//ClientConn rpc 连接
func (node *NodeEnv) ClientConn() error {
	node.Lock()
	defer node.Unlock()
	if node.Connect != nil {
		return nil
	}

	if node.connFail != nil {
		return node.connFail
	}
	err := node.verify()
	if err != nil {
		return err
	}
	//response to do submit
	if node.waitConn == nil {
		node.waitConn = sync.NewCond(node)
		go func() {
			// conn := &ClientConn{nil, true}
			// err := conn.Dial(c.endpointConf)
			conn, err := node.grpcConnection()
			if err != nil {
				err = errors.WithMessage(err, fmt.Sprintf("client failed to connect to %s", node.Address))
			}
			node.Lock()
			if err != nil {
				node.connFail = err
			} else {
				node.Connect = conn
			}
			node.Unlock()

			node.waitConn.Broadcast()
			fmt.Println("connect node addr:", node.Address)

			// if err != nil {
			// 	node.reset(ctx)
			// }
		}()

		node.waitConn.Wait()
		node.waitConn = nil

	} else {
		node.waitConn.Wait()
	}

	if node.Connect != nil {
		return nil
	} else {
		return node.connFail
	}
	// conn, err := node.grpcConnection()
	// if err != nil {
	// 	return errors.WithMessage(err, fmt.Sprintf("endorser client failed to connect to %s", node.Address))
	// }
	// node.Connect = conn
	// return nil
}

//grpc 关闭连接
func (node *NodeEnv) CloseConn() {
	// node.Lock()
	// defer node.Unlock()
	node.Lock()
	defer node.Unlock()

	if node.waitConn != nil {
		node.waitConn.Wait()
	}
	if node.Connect != nil {
		node.Connect.Close()
		node.Connect = nil
	}
}

//
func (node *NodeEnv) grpcConnection() (*grpc.ClientConn, error) {
	logger.Debug("创建grpc 连接")
	var dialOpts []grpc.DialOption
	if node.TLS {
		tls, err := credentials.NewClientTLSFromFile(node.RootCertFile, node.HostnameOverride)
		if err != nil {
			return nil, err
		}
		fmt.Println("tls true")
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(tls), grpc.WithBlock())

	} else {
		fmt.Println("tls false")
		dialOpts = append(dialOpts, grpc.WithInsecure(), grpc.WithBlock())
		// dialOpts = append(dialOpts, grpc.WithTransportCredentials(tls), grpc.WithBlock())
	}
	// tls, err := credentials.NewClientTLSFromFile(node.RootCertFile, node.HostnameOverride)
	// if err != nil {
	// 	return nil, err
	// }
	// dialOpts = append(dialOpts, grpc.WithTransportCredentials(tls), grpc.WithBlock())

	ctx, cancel := context.WithTimeout(context.Background(), node.ConnTimeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, node.Address, dialOpts...)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create new connection")
	}
	return conn, nil
}

//VerifyConn 校验状态
func (node *NodeEnv) VerifyConn() error {

	if node.Connect == nil {
		return fmt.Errorf("Conn not inited")
	}

	s := node.Connect.GetState()

	if s != grpc_conn.Ready {
		return fmt.Errorf("Conn is not ready: <%s>", s)
	}

	return nil
}

func (node *NodeEnv) reset(ctx context.Context) {

	to := node.resetInterval
	if int64(to) == 0 {
		to = defaultReconnectInterval
	}

	select {
	case <-time.After(to):
		node.Lock()
		node.connFail = nil
		node.Unlock()
	case <-ctx.Done():
	}
}
