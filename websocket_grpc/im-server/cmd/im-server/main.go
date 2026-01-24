package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/bootstrap"
	"go.uber.org/zap"
)

func main() {
	cfg := bootstrap.DefaultConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app, err := bootstrap.New(ctx, cfg)
	if err != nil {
		panic(err)
	}

	srvErr := make(chan error, 1)
	srvDone := make(chan struct{})
	go func() {
		defer close(srvDone)
		if err := app.HTTP.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		app.Log.Info("shutdown signal received")
	case err := <-srvErr:
		app.Log.Error("http server error", zap.Error(err))
		stop() // 触发退出流程
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = app.HTTP.Stop(shutdownCtx)
	<-srvDone
	_ = app.Log.Sync()
	os.Exit(0)
}
