package registry

import (
	"context"
	"sync"

	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/app/transport/ws/contract"
)

type registry struct {
	mu sync.Mutex
	// user -> device -> session
	sessions map[string]map[string]contract.Session
}

func NewRegistry() contract.Registry {
	return &registry{
		sessions: make(map[string]map[string]contract.Session),
	}
}

func (r *registry) Bind(ctx context.Context, session contract.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.sessions[session.UserID()]; !ok {
		r.sessions[session.UserID()] = make(map[string]contract.Session)
	}
	r.sessions[session.UserID()][session.DeviceID()] = session
	return nil
}

func (r *registry) Unbind(ctx context.Context, userID, DeviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	devices, ok := r.sessions[userID]
	if !ok {
		return nil
	}

	delete(devices, DeviceID)
	if len(devices) == 0 {
		delete(r.sessions, userID)
	}

	return nil
}

func (r *registry) GetUserSessions(ctx context.Context, userID string) ([]contract.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	sessions := make([]contract.Session, 0, len(r.sessions[userID]))
	for _, session := range r.sessions[userID] {
		sessions = append(sessions, session)
	}

	return sessions, nil
}
