package handler

import (
	"context"
	"fmt"

	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws/contract"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws/protocol"
	imv1 "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/gen/im/v1"
)

type GroupMessageHandler struct {
	reg   contract.Registry
	store *GroupStore
}

func NewGroupMessageHandler(reg contract.Registry, store *GroupStore) *GroupMessageHandler {
	return &GroupMessageHandler{reg: reg, store: store}
}

func (h *GroupMessageHandler) HandleMessage(ctx context.Context, sess contract.Session, msg *imv1.ClientEnvelope) error {
	p := msg.GetGroupMessage()
	if p == nil {
		return fmt.Errorf("missing group_message payload")
	}
	if p.GetUuid() == "" {
		return fmt.Errorf("uuid is required")
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

	if h.store == nil {
		// No membership resolver wired yet.
		return nil
	}

	memberIDs := h.store.ListMembers(p.GetUuid())
	if len(memberIDs) == 0 {
		return nil
	}

	// Deliver to each member's online sessions.
	deliver := &imv1.ServerEnvelope{
		TraceId: msg.GetTraceId(),
		Payload: &imv1.ServerEnvelope_DeliverGroupMessage{
			DeliverGroupMessage: &imv1.DeliverGroupMessage{
				GroupUuid: p.GetUuid(),
				From:      sess.UserID(),
				Message:   p.GetMessage(),
			},
		},
	}
	deliverBytes, err := protocol.EncodeServerMessage(deliver)
	if err != nil {
		return err
	}

	for _, uid := range memberIDs {
		targetSessions, err := h.reg.GetUserSessions(ctx, uid)
		if err != nil {
			return err
		}
		for _, ts := range targetSessions {
			if ts == nil {
				continue
			}
			_ = ts.Send(deliverBytes)
		}
	}
	return nil
}

