package handlers

import (
	"context"

	v1 "github.com/gussf/poc-grpc-gateway/proto/gen/go"
)

type Echoer struct {
	v1.UnimplementedEchoerServer
}

func NewEchoerHandler() Echoer {
	return Echoer{}
}

func (e Echoer) PostEcho(ctx context.Context, msg *v1.StringMessage) (*v1.StringMessage, error) {
	return &v1.StringMessage{
		Value: msg.GetValue(),
	}, nil
}

func (e Echoer) GetEcho(ctx context.Context, msg *v1.StringMessage) (*v1.StringMessage, error) {
	return &v1.StringMessage{
		Value: msg.GetValue(),
	}, nil
}
