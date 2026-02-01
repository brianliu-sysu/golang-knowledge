package handler

import (
	"context"

	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws/contract"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws/protocol"
	imv1 "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/gen/im/v1"
)

type EchoHandler struct {
}

func (h *EchoHandler) HandleMessage(ctx context.Context, sess contract.Session, msg *imv1.ClientEnvelope) error {
	resp := &imv1.ServerEnvelope{
		TraceId: msg.TraceId,
		Payload: &imv1.ServerEnvelope_Echo{
			Echo: &imv1.Echo{Message: msg.Payload.(*imv1.ClientEnvelope_Echo).Echo.Message},
		},
	}
	out, err := protocol.EncodeServerMessage(resp)
	if err != nil {
		return err
	}
	return sess.Send(out)
}
