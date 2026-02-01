package inbox

import "time"

type Inbox struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	MessageID int64     `json:"message_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
