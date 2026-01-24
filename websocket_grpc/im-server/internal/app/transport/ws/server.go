package ws

import (
	"net/http"
	"time"

	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/services/auth"
	"go.uber.org/zap"
)

type ServerOptions struct {
	NodeID         string
	Path           string
	ReadLimitBytes int64
	SendQueueSize  int
	WriteTimeout   time.Duration
	PingInterval   time.Duration
	PongWait       time.Duration
	Logger         *zap.Logger
	Registry       Registry
	Authenticator  auth.Authenticator
}

type Server struct {
	h *Handler
}

func NewServer(opts ServerOptions) *Server {
	return &Server{
		h: NewHandler(HandlerOptions{
			NodeID:         opts.NodeID,
			ReadLimitBytes: int(opts.ReadLimitBytes),
			SendQueueSize:  opts.SendQueueSize,
			WriteTimeout:   opts.WriteTimeout,
			PingInterval:   opts.PingInterval,
			PongWait:       opts.PongWait,
			Logger:         opts.Logger,
			Registry:       opts.Registry,
			Authenticator:  opts.Authenticator,
		}),
	}
}

// ServeHTTP makes Server implement http.Handler so it can be mounted directly on an http mux.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.h.ServeHTTP(w, r)
}
