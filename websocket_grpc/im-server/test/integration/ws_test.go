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
		Type:    imv1.MessageType_ECHO,
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

func TestWS_SingleMessage_OK(t *testing.T) {
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

	// receiver u2
	h2 := make(http.Header)
	h2.Add("Authorization", "Bearer test:u2:d1")
	conn2, _, err := websocket.DefaultDialer.Dial(u, h2)
	if err != nil {
		t.Fatalf("failed to dial websocket (u2): %v", err)
	}
	defer conn2.Close()

	// sender u1
	h1 := make(http.Header)
	h1.Add("Authorization", "Bearer test:u1:d1")
	conn1, _, err := websocket.DefaultDialer.Dial(u, h1)
	if err != nil {
		t.Fatalf("failed to dial websocket (u1): %v", err)
	}
	defer conn1.Close()

	req := &imv1.ClientEnvelope{
		TraceId: "t-single-1",
		Type:    imv1.MessageType_SINGLE_MESSAGE,
		Payload: &imv1.ClientEnvelope_SingleMessage{
			SingleMessage: &imv1.SingleMessage{
				To:      "u2",
				Message: []byte("hi u2"),
			},
		},
	}
	reqBytes, err := protocol.EncodeClientMessage(req)
	if err != nil {
		t.Fatalf("failed to encode client message: %v", err)
	}
	if err := conn1.WriteMessage(websocket.BinaryMessage, reqBytes); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}

	// sender gets ack
	_, ackBytes, err := conn1.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read ack: %v", err)
	}
	var ackEnv imv1.ServerEnvelope
	if err := proto.Unmarshal(ackBytes, &ackEnv); err != nil {
		t.Fatalf("failed to unmarshal ack: %v", err)
	}
	if _, ok := ackEnv.Payload.(*imv1.ServerEnvelope_AckResp); !ok {
		t.Fatalf("expected ack_resp, got %T", ackEnv.Payload)
	}

	// receiver gets delivery
	_, delBytes, err := conn2.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read delivery: %v", err)
	}
	var delEnv imv1.ServerEnvelope
	if err := proto.Unmarshal(delBytes, &delEnv); err != nil {
		t.Fatalf("failed to unmarshal delivery: %v", err)
	}
	ds, ok := delEnv.Payload.(*imv1.ServerEnvelope_DeliverSingleMessage)
	if !ok {
		t.Fatalf("expected deliver_single_message, got %T", delEnv.Payload)
	}
	if ds.DeliverSingleMessage.GetFrom() != "u1" {
		t.Fatalf("expected from=u1, got %s", ds.DeliverSingleMessage.GetFrom())
	}
	if !bytes.Equal(ds.DeliverSingleMessage.GetMessage(), []byte("hi u2")) {
		t.Fatalf("unexpected message: %s", string(ds.DeliverSingleMessage.GetMessage()))
	}
}
