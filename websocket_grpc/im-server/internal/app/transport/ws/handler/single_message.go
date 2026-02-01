package handler

import (
	"context"
	"fmt"

	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws/contract"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws/protocol"
	imv1 "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/gen/im/v1"
)

type SingleMessageHandler struct {
	reg contract.Registry
}

func NewSingleMessageHandler(reg contract.Registry) *SingleMessageHandler {
	return &SingleMessageHandler{reg: reg}
}

func (h *SingleMessageHandler) HandleMessage(ctx context.Context, sess contract.Session, msg *imv1.ClientEnvelope) error {
	p := msg.GetSingleMessage()
	if p == nil {
		return fmt.Errorf("missing single_message payload")
	}
	if p.GetTo() == "" {
		return fmt.Errorf("to is required")
	}
	if h.reg == nil {
		return fmt.Errorf("registry is nil")
	}

	// Ack to sender.
	ack := &imv1.ServerEnvelope{
		TraceId: msg.GetTraceId(),
		Payload: &imv1.ServerEnvelope_AckResp{
			AckResp: &imv1.AckResp{Status: 0, Message: "ok"},
		},
	}
	ackBytes, err := protocol.EncodeServerMessage(ack)
	if err != nil {
		return err
	}
	if err := sess.Send(ackBytes); err != nil {
		return err
	}

	// Deliver to recipient's online sessions.
	deliver := &imv1.ServerEnvelope{
		TraceId: msg.GetTraceId(),
		Payload: &imv1.ServerEnvelope_DeliverSingleMessage{
			DeliverSingleMessage: &imv1.DeliverSingleMessage{
				From:    sess.UserID(),
				Message: p.GetMessage(),
			},
		},
	}
	deliverBytes, err := protocol.EncodeServerMessage(deliver)
	if err != nil {
		return err
	}

	targetSessions, err := h.reg.GetUserSessions(ctx, p.GetTo())
	if err != nil {
		return err
	}
	for _, ts := range targetSessions {
		if ts == nil {
			continue
		}
		_ = ts.Send(deliverBytes)
	}
	return nil
}

