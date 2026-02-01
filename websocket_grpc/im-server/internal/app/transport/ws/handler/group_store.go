package handler

import "sync"

// GroupStore is a minimal in-memory group membership store.
// It's intentionally simple for now; replace with repository-based implementation later.
type GroupStore struct {
	mu      sync.RWMutex
	members map[string]map[string]struct{} // groupUUID -> userID set
}

func NewGroupStore() *GroupStore {
	return &GroupStore{members: make(map[string]map[string]struct{})}
}

func (s *GroupStore) AddMember(groupUUID, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.members[groupUUID]; !ok {
		s.members[groupUUID] = make(map[string]struct{})
	}
	s.members[groupUUID][userID] = struct{}{}
}

func (s *GroupStore) ListMembers(groupUUID string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m := s.members[groupUUID]
	if len(m) == 0 {
		return nil
	}
	out := make([]string, 0, len(m))
	for uid := range m {
		out = append(out, uid)
	}
	return out
}

