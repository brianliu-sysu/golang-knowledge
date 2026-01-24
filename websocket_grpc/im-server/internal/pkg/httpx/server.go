package httpx

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Options struct {
	Addr    string
	Handler http.Handler
	Logger  *zap.Logger
}

type Server struct {
	srv *http.Server
	log *zap.Logger
}

func NewServer(opts Options) *Server {
	srv := &http.Server{
		Addr:         opts.Addr,
		Handler:      opts.Handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		srv: srv,
		log: opts.Logger,
	}
}

func (s *Server) Start() error {
	s.log.Info("starting HTTP server", zap.String("addr", s.srv.Addr))
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("stopping HTTP server")
	return s.srv.Shutdown(ctx)
}
