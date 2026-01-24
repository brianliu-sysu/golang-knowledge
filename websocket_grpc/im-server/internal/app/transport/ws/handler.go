package ws

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws/protocol"
	imv1 "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/gen/im/v1"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/pkg/observability"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/services/auth"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type HandlerOptions struct {
	NodeID         string
	ReadLimitBytes int
	SendQueueSize  int
	WriteTimeout   time.Duration
	PingInterval   time.Duration
	PongWait       time.Duration
	Logger         *zap.Logger
	Registry       Registry
	Authenticator  auth.Authenticator
}

type Handler struct {
	options  HandlerOptions
	upgrader websocket.Upgrader
}

func NewHandler(options HandlerOptions) *Handler {
	return &Handler{
		options: options,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  options.ReadLimitBytes,
			WriteBufferSize: options.ReadLimitBytes,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func bearerTokenFromRequest(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.options.Logger.Info("Handling WebSocket connection")
	ctx := r.Context()
	token := bearerTokenFromRequest(r)
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.options.Authenticator.Authenticate(ctx, token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	observability.WSConnections.Inc()
	conn.SetReadLimit(int64(h.options.ReadLimitBytes))
	_ = conn.SetReadDeadline(time.Now().Add(h.options.PongWait))
	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(h.options.PongWait))
		return nil
	})
	s := NewSession(user.UserID, user.DeviceID, h.options.NodeID, conn, h.options.ReadLimitBytes, h.options.WriteTimeout)
	h.options.Registry.Bind(ctx, s)

	connCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer cancel()
		h.writeLoop(connCtx, s)
	}()
	go func() {
		defer wg.Done()
		defer cancel()
		h.readLoop(connCtx, s)
	}()
	wg.Wait()
	_ = s.Close()
	h.options.Registry.Unbind(ctx, user.UserID, user.DeviceID)
	observability.WSConnections.Dec()
}

func (h *Handler) writeLoop(ctx context.Context, s *Session) {
	ticker := time.NewTicker(h.options.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = s.conn.SetWriteDeadline(time.Now().Add(h.options.WriteTimeout))
			if err := s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				h.options.Logger.Info("failed to write ping", zap.Error(err))
				return
			}
		case <-ctx.Done():
			return
		case payload, ok := <-s.send:
			if !ok {
				return
			}

			_ = s.conn.SetWriteDeadline(time.Now().Add(h.options.WriteTimeout))
			if err := s.conn.WriteMessage(websocket.BinaryMessage, payload); err != nil {
				h.options.Logger.Info("failed to write message", zap.Error(err))
				return
			}
		}
	}
}

func (h *Handler) readLoop(ctx context.Context, s *Session) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			mt, payload, err := s.conn.ReadMessage()
			if err != nil {
				h.options.Logger.Info("failed to read message", zap.Error(err))
				return
			}
			if mt != websocket.BinaryMessage {
				h.options.Logger.Info("receive unexpect data", zap.Any("type", mt), zap.Any("data", payload))
				_ = h.SendError(ctx, s, "", "BAD_REQUEST", "Invalid message type")
				continue
			}

			env, err := protocol.DecodeClientMessage(payload)
			if err != nil {
				observability.WSBadProto.Inc()
				_ = h.SendError(ctx, s, "", "BAD_PROTO", "Protocol error")
				continue
			}
			// MVP：只支持 Echo
			switch p := env.Payload.(type) {
			case *imv1.ClientEnvelope_Echo:
				resp := &imv1.ServerEnvelope{
					TraceId: env.TraceId,
					Payload: &imv1.ServerEnvelope_Echo{
						Echo: &imv1.Echo{Message: p.Echo.Message},
					},
				}
				out, e := protocol.EncodeServerMessage(resp)
				if e != nil {
					_ = h.SendError(ctx, s, env.TraceId, "INTERNAL", "encode failed")
					continue
				}
				_ = s.Send(out)
			default:
				_ = h.SendError(ctx, s, env.TraceId, "UNSUPPORTED", "payload not supported in milestone 1")
			}
		}
	}
}

func (h *Handler) SendError(ctx context.Context, s *Session, traceID, code, message string) error {
	resp := &imv1.ServerEnvelope{
		TraceId: traceID,
		Payload: &imv1.ServerEnvelope_Error{
			Error: &imv1.Error{
				Code:    code,
				Message: message,
			},
		},
	}

	payload, err := protocol.EncodeServerMessage(resp)
	if err != nil {
		h.options.Logger.Info("failed to encode error message", zap.Error(err))
		return err
	}
	return s.Send(payload)
}
