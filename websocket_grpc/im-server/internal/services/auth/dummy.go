package auth

import (
	"context"
	"errors"
	"strings"
)

var ErrUnauthorized = errors.New("unauthorized")

type DummyAuthenticator struct{}

func NewDummyAuthenticator() Authenticator {
	return &DummyAuthenticator{}
}

func (a *DummyAuthenticator) Authenticate(ctx context.Context, token string) (*Principal, error) {
	if token == "" {
		return nil, ErrUnauthorized
	}

	parts := strings.Split(token, ":")
	if len(parts) != 3 || parts[0] != "test" {
		return nil, ErrUnauthorized
	}

	return &Principal{
		UserID:   parts[1],
		DeviceID: parts[2],
	}, nil
}
