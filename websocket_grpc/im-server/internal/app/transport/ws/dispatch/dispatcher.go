package dispatch

import (
	"context"
	"fmt"

	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws/contract"
	wshandler "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws/handler"
	imv1 "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/gen/im/v1"
)

// NewDefaultDispatcher builds a dispatcher using per-app injected dependencies (e.g. registry),
// avoiding global singletons that could drift from the wiring layer.
func NewDefaultDispatcher(reg contract.Registry) Dispatcher {
	store := wshandler.NewGroupStore()
	d := &dispatcher{handlers: make(map[imv1.MessageType]MessageHandler)}
	d.RegisterHandler(imv1.MessageType_ECHO, &wshandler.EchoHandler{})
	d.RegisterHandler(imv1.MessageType_CREATE_GROUP, wshandler.NewGroupHandler(reg, store))
	d.RegisterHandler(imv1.MessageType_LIST_GROUPS, wshandler.NewGroupHandler(reg, store))
	d.RegisterHandler(imv1.MessageType_SINGLE_MESSAGE, wshandler.NewSingleMessageHandler(reg))
	d.RegisterHandler(imv1.MessageType_GROUP_MESSAGE, wshandler.NewGroupMessageHandler(reg, store))
	return d
}

type Dispatcher interface {
	Dispatch(ctx context.Context, sess contract.Session, msg *imv1.ClientEnvelope) error
}

type MessageHandler interface {
	HandleMessage(ctx context.Context, sess contract.Session, msg *imv1.ClientEnvelope) error
}

type dispatcher struct {
	handlers map[imv1.MessageType]MessageHandler
}

func NewDispatcher() Dispatcher {
	return &dispatcher{handlers: make(map[imv1.MessageType]MessageHandler)}
}

func NewDispatcherWithHandlers(handlers map[imv1.MessageType]MessageHandler) Dispatcher {
	return &dispatcher{handlers: handlers}
}

func (d *dispatcher) RegisterHandler(messageType imv1.MessageType, handler MessageHandler) {
	d.handlers[messageType] = handler
}

func (d *dispatcher) Dispatch(ctx context.Context, sess contract.Session, msg *imv1.ClientEnvelope) error {
	// Prefer explicit name-based routing (works well across languages),
	// but fall back to payload type when name is empty.
	if msg.GetType() != imv1.MessageType_UNSPECIFIED {
		h, ok := d.handlers[msg.GetType()]
		if !ok {
			return fmt.Errorf("handler not found for type: %v", msg.GetType())
		}
		return h.HandleMessage(ctx, sess, msg)
	}

	// Fallback routing by payload when type is unset.
	switch msg.GetPayload().(type) {
	case *imv1.ClientEnvelope_Echo:
		h, ok := d.handlers[imv1.MessageType_ECHO]
		if !ok {
			return fmt.Errorf("handler not found for payload: echo")
		}
		return h.HandleMessage(ctx, sess, msg)
	default:
		return fmt.Errorf("unsupported payload")
	}
}
