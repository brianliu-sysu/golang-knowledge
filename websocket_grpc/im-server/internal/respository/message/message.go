package message

import (
	"context"

	domainmessage "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/domain/message"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *domainmessage.Message) error
	GetMessage(ctx context.Context, id int64) (*domainmessage.Message, error)
	UpdateMessage(ctx context.Context, message *domainmessage.Message) error
	DeleteMessage(ctx context.Context, id int64) error
	ListMessages(ctx context.Context, userID int64, page int, pageSize int) ([]*domainmessage.Message, error)
	ListMessagesByGroupID(ctx context.Context, groupID int64, page int, pageSize int) ([]*domainmessage.Message, error)
}
