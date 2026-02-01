package group

import (
	"context"

	domaingroup "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/domain/group"
)

type GroupRepository interface {
	CreateGroup(ctx context.Context, group *domaingroup.Group) error
	GetGroup(ctx context.Context, id int64) (*domaingroup.Group, error)
	UpdateGroup(ctx context.Context, group *domaingroup.Group) error
	DeleteGroup(ctx context.Context, id int64) error
	ListGroups(ctx context.Context, page int, pageSize int) ([]*domaingroup.Group, error)
	ListGroupMembers(ctx context.Context, groupID int64, page int, pageSize int) ([]*domaingroup.GroupMember, error)
}
