package contract

import "context"

// Session is the minimal, framework-agnostic view of a websocket connection
// that business logic may need.
type Session interface {
	UserID() string
	DeviceID() string
	NodeID() string
	Send(data []byte) error
}

// Registry is the minimal interface for looking up online sessions.
// Keep it small to avoid pulling ws implementation details into business code.
type Registry interface {
	Bind(ctx context.Context, session Session) error
	Unbind(ctx context.Context, userID, deviceID string) error
	GetUserSessions(ctx context.Context, userID string) ([]Session, error)
}

