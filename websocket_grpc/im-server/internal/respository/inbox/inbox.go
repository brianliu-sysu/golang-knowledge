package inbox

import (
	"context"

	domaininbox "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/domain/inbox"
)

type InboxRepository interface {
	CreateInbox(ctx context.Context, inbox *domaininbox.Inbox) error
	GetInbox(ctx context.Context, id int64) (*domaininbox.Inbox, error)
	UpdateInbox(ctx context.Context, inbox *domaininbox.Inbox) error
	DeleteInbox(ctx context.Context, id int64) error
	ListInboxes(ctx context.Context, userID int64, page int, pageSize int) ([]*domaininbox.Inbox, error)
}
