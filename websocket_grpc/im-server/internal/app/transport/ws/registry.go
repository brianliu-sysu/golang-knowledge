package ws

import (
	"context"
	"sync"
)

type Registry interface {
	Bind(ctx context.Context, session *Session) error
	Unbind(ctx context.Context, userID, DeviceID string) error
	GetUserSessions(ctx context.Context, userID string) ([]*Session, error)
}

type registry struct {
	mu sync.Mutex
	// user -> device -> session
	sessions map[string]map[string]*Session
}

func NewRegistry() Registry {
	return &registry{
		sessions: make(map[string]map[string]*Session),
	}
}

func (r *registry) Bind(ctx context.Context, session *Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.sessions[session.UserID()]; !ok {
		r.sessions[session.UserID()] = make(map[string]*Session)
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

func (r *registry) GetUserSessions(ctx context.Context, userID string) ([]*Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	sessions := make([]*Session, 0, len(r.sessions[userID]))
	for _, session := range r.sessions[userID] {
		sessions = append(sessions, session)
	}

	return sessions, nil
}
