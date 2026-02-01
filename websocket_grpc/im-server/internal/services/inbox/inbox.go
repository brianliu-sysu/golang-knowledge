package inbox

import (
	"context"

	domaininbox "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/domain/inbox"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/respository/inbox"
)

type InboxService struct {
	inboxRepository inbox.InboxRepository
}

func NewInboxService(inboxRepository inbox.InboxRepository) *InboxService {
	return &InboxService{inboxRepository: inboxRepository}
}

func (s *InboxService) CreateInbox(ctx context.Context, inbox *domaininbox.Inbox) error {
	return s.inboxRepository.CreateInbox(ctx, inbox)
}