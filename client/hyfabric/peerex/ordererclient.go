package peerex

import (
	"context"

	cb "github.com/hyperledger/fabric/protos/common"
	ab "github.com/hyperledger/fabric/protos/orderer"
	"github.com/pkg/errors"
)

type BroadcastClient interface {
	//Send data to orderer
	Send(env *cb.Envelope) error
	Close() error
}

type broadcastClient struct {
	client ab.AtomicBroadcast_BroadcastClient
}

func (order *OrderEnv) NewBroadcastClient() (BroadcastClient, error) {

	bc, err := ab.NewAtomicBroadcastClient(order.Connect).Broadcast(context.TODO())
	if err != nil {
		return nil, err
	}

	return &broadcastClient{client: bc}, nil
}

//Send data to orderer
func (s *broadcastClient) Send(env *cb.Envelope) error {
	if err := s.client.Send(env); err != nil {
		return errors.WithMessage(err, "could not send")
	}

	msg, err := s.client.Recv()
	if err != nil {
		return err
	}
	if msg.Status != cb.Status_SUCCESS {
		return errors.Errorf("got unexpected status: %v -- %s", msg.Status, msg.Info)
	}
	return nil
}

func (s *broadcastClient) Close() error {
	return s.client.CloseSend()
}
