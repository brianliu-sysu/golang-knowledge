package integration

import (
	"bytes"
	"testing"

	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/bootstrap"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws/protocol"
	imv1 "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/gen/im/v1"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/services/auth"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	proto "google.golang.org/protobuf/proto"
)

func TestWS_Unauthenticated(t *testing.T) {
	cfg := bootstrap.DefaultConfig()
	reg := ws.NewRegistry()
	authn := auth.NewDummyAuthenticator()
	mux := http.NewServeMux()
	srv := ws.NewServer(ws.ServerOptions{
		NodeID:         cfg.NodeID,
		Path:           cfg.WSPath,
		ReadLimitBytes: cfg.ReadLimitBytes,
		SendQueueSize:  cfg.SendQueueSize,
		WriteTimeout:   cfg.WriteTimeout,
		PingInterval:   cfg.PingInterval,
		PongWait:       cfg.PongWait,
		Logger:         zap.NewNop(),
		Registry:       reg,
		Authenticator:  authn,
	})
	mux.Handle(cfg.WSPath, srv)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	u := "ws://" + strings.TrimPrefix(ts.URL, "http://") + cfg.WSPath
	_, resp, err := websocket.DefaultDialer.Dial(u, nil)
	if err == nil {
		t.Fatalf("should be error, got %v", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("should be unauthorized, got %v", resp.StatusCode)
	}
}

func TestWS_Echo_OK(t *testing.T) {
	cfg := bootstrap.DefaultConfig()
	reg := ws.NewRegistry()
	authn := auth.NewDummyAuthenticator()
	mux := http.NewServeMux()
	srv := ws.NewServer(ws.ServerOptions{
		NodeID:         cfg.NodeID,
		Path:           cfg.WSPath,
		ReadLimitBytes: cfg.ReadLimitBytes,
		SendQueueSize:  cfg.SendQueueSize,
		WriteTimeout:   cfg.WriteTimeout,
		PingInterval:   cfg.PingInterval,
		PongWait:       cfg.PongWait,
		Logger:         zap.NewNop(),
		Registry:       reg,
		Authenticator:  authn,
	})
	mux.Handle(cfg.WSPath, srv)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	u := "ws://" + strings.TrimPrefix(ts.URL, "http://") + cfg.WSPath
	h := make(http.Header)
	h.Add("Authorization", "Bearer test:u1:d1")
	conn, _, err := websocket.DefaultDialer.Dial(u, h)
	if err != nil {
		t.Fatalf("failed to dial websocket: %v", err)
	}
	defer conn.Close()

	req := &imv1.ClientEnvelope{
		TraceId: "123",
		Payload: &imv1.ClientEnvelope_Echo{
			Echo: &imv1.Echo{
				Message: []byte("hello"),
			},
		},
	}
	reqBytes, err := protocol.EncodeClientMessage(req)
	if err != nil {
		t.Fatalf("failed to encode client message: %v", err)
	}
	err = conn.WriteMessage(websocket.BinaryMessage, reqBytes)
	if err != nil {
		t.Fatalf("failed to write message: %v", err)
	}

	mt, resp, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}
	if mt != websocket.BinaryMessage {
		t.Fatalf("should be binary message, got %v", mt)
	}
	var respMsg imv1.ServerEnvelope
	err = proto.Unmarshal(resp, &respMsg)
	if err != nil {
		t.Fatalf("failed to unmarshal server message: %v", err)
	}
	if respMsg.TraceId != "123" {
		t.Fatalf("should be 123, got %v", respMsg.TraceId)
	}

	pe, ok := respMsg.Payload.(*imv1.ServerEnvelope_Echo)
	if !ok {
		t.Fatalf("should be echo, got %v", respMsg.Payload)
	}
	if !bytes.Equal(pe.Echo.Message, []byte("hello")) {
		t.Fatalf("should be hello, got %v", pe.Echo.Message)
	}
}
