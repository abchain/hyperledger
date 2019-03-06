package client

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"hyperledger.abchain.org/core/utils"
	"time"
)

type ClientConn struct {
	C         *grpc.ClientConn
	BlockConn bool
}

const defaultTimeout = time.Second * 3

func (conn *ClientConn) Dial(conf map[string]string) error {

	addr := conf["server"]
	if addr == "" {
		return fmt.Errorf("Server is not specified")
	}

	var opts []grpc.DialOption

	tls := conf["tlsenabled"]
	if tls == "true" {
		cert := conf["certfile"]
		hostName := conf["hostname"]

		creds, err := credentials.NewClientTLSFromFile(utils.CanonicalizeFilePath(cert), hostName)
		if err != nil {
			return fmt.Errorf("Failed to create TLS credentials %v", err)
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	//Todo: set timeout?
	opts = append(opts, grpc.WithTimeout(defaultTimeout))
	if conn.BlockConn {
		opts = append(opts, grpc.WithBlock())
	}
	c, err := grpc.Dial(addr, opts...)
	if err != nil {
		return fmt.Errorf("Dial fail: %s", err)
	}

	conn.C = c
	return nil
}
