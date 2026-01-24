package auth

import "context"

type Principal struct {
	UserID   string `json:"user_id"`
	DeviceID string `json:"device_id"`
}

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (*Principal, error)
}
