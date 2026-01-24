package bootstrap

import (
	"context"
	"net/http"

	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/pkg/httpx"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/pkg/observability"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/services/auth"

	"go.uber.org/zap"
)

type App struct {
	HTTP *httpx.Server
	WS   *ws.Server
	Log  *zap.Logger
}

func New(ctx context.Context, cfg *Config) (*App, error) {
	log := observability.NewLogger()

	reg := ws.NewRegistry()
	authn := auth.NewDummyAuthenticator()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.Handle("/metrics", observability.MetricsHandler())

	wsServer := ws.NewServer(ws.ServerOptions{
		NodeID:         cfg.NodeID,
		Path:           cfg.WSPath,
		ReadLimitBytes: cfg.ReadLimitBytes,
		SendQueueSize:  cfg.SendQueueSize,
		WriteTimeout:   cfg.WriteTimeout,
		PingInterval:   cfg.PingInterval,
		PongWait:       cfg.PongWait,
		Logger:         log,
		Registry:       reg,
		Authenticator:  authn,
	})
	mux.Handle(cfg.WSPath, wsServer)

	httpServer := httpx.NewServer(httpx.Options{
		Addr:    cfg.HTTPAddr,
		Handler: mux,
		Logger:  log,
	})

	return &App{HTTP: httpServer, WS: wsServer, Log: log}, nil
}
