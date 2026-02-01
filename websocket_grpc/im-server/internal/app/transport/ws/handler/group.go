package handler

import (
	"context"
	"fmt"

	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws/contract"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws/protocol"
	imv1 "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/gen/im/v1"
)

type GroupHandler struct {
	registry contract.Registry
	store    *GroupStore
}

func NewGroupHandler(registry contract.Registry, store *GroupStore) *GroupHandler {
	return &GroupHandler{registry: registry, store: store}
}

func (h *GroupHandler) HandleMessage(ctx context.Context, sess contract.Session, msg *imv1.ClientEnvelope) error {
	switch msg.GetType() {
	case imv1.MessageType_CREATE_GROUP:
		return h.createGroup(ctx, sess, msg)
	case imv1.MessageType_LIST_GROUPS:
		return h.listGroup(ctx, sess, msg)
	default:
		return fmt.Errorf("unsupported message type: %v", msg.GetType())
	}
}

func (h *GroupHandler) createGroup(ctx context.Context, sess contract.Session, msg *imv1.ClientEnvelope) error {
	p := msg.GetCreateGroup()
	if p == nil {
		return fmt.Errorf("missing create_group payload")
	}
	if p.GetUuid() == "" {
		return fmt.Errorf("uuid is required")
	}

	// MVP: in-memory membership only.
	if h.store != nil {
		h.store.AddMember(p.GetUuid(), sess.UserID())
	}

	resp := &imv1.ServerEnvelope{
		TraceId: msg.GetTraceId(),
		Payload: &imv1.ServerEnvelope_AckResp{
			AckResp: &imv1.AckResp{Status: 0, Message: "ok"},
		},
	}
	out, err := protocol.EncodeServerMessage(resp)
	if err != nil {
		return err
	}
	return sess.Send(out)
}

func (h *GroupHandler) listGroup(ctx context.Context, sess contract.Session, msg *imv1.ClientEnvelope) error {
	// TODO: wire real group repository and return ListGroupsResp.
	resp := &imv1.ServerEnvelope{
		TraceId: msg.GetTraceId(),
		Payload: &imv1.ServerEnvelope_AckResp{
			AckResp: &imv1.AckResp{Status: 0, Message: "ok"},
		},
	}
	out, err := protocol.EncodeServerMessage(resp)
	if err != nil {
		return err
	}
	return sess.Send(out)
}
