package message

import "time"

// Message is the core domain model for chat messages.
// Keep this minimal for now; expand as business requirements grow.
type Message struct {
	ID        int64
	FromUser  int64
	ToUser    int64
	GroupID   int64
	Content   string
	CreatedAt time.Time
}

