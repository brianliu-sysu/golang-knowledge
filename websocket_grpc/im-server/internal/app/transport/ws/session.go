package ws

import (
	"errors"
	"time"

	"github.com/gorilla/websocket"
)

var ErrBackPressure = errors.New("backpressure")

type Session struct {
	userID       string
	deviceID     string
	nodeID       string
	conn         *websocket.Conn
	send         chan []byte
	WriteTimeout time.Duration
}

func NewSession(userID, deviceID, nodeID string,
	conn *websocket.Conn, sendQueueSize int, writeTimeout time.Duration) *Session {
	return &Session{
		userID:       userID,
		deviceID:     deviceID,
		nodeID:       nodeID,
		conn:         conn,
		send:         make(chan []byte, sendQueueSize),
		WriteTimeout: writeTimeout,
	}
}

func (s *Session) Send(data []byte) error {
	select {
	case s.send <- data:
		return nil
	default:
		return ErrBackPressure
	}
}

func (s *Session) Close() error {
	close(s.send)
	return s.conn.Close()
}

func (s *Session) UserID() string {
	return s.userID
}

func (s *Session) DeviceID() string {
	return s.deviceID
}

func (s *Session) NodeID() string {
	return s.nodeID
}
