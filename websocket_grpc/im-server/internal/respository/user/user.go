package user

import (
	"context"

	domainuser "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/domain/user"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domainuser.User) error
	GetUser(ctx context.Context, id int64) (*domainuser.User, error)
	UpdateUser(ctx context.Context, user *domainuser.User) error
	DeleteUser(ctx context.Context, id int64) error
}
