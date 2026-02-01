package postgres

import (
	"context"

	domaininbox "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/domain/inbox"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/respository/inbox"
	"github.com/jackc/pgx/v5/pgxpool"
)

type inboxRepository struct {
	pool *pgxpool.Pool
}

func NewInboxRepository(pool *pgxpool.Pool) inbox.InboxRepository {
	return &inboxRepository{pool: pool}
}

func (r *inboxRepository) CreateInbox(ctx context.Context, inbox *domaininbox.Inbox) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO inboxes (id, user_id, message_id, created_at) VALUES ($1, $2, $3, $4)", inbox.ID, inbox.UserID, inbox.MessageID, inbox.CreatedAt)
	return err
}

func (r *inboxRepository) GetInbox(ctx context.Context, id int64) (*domaininbox.Inbox, error) {
	row := r.pool.QueryRow(ctx, "SELECT id, user_id, message_id, created_at FROM inboxes WHERE id = $1", id)
	var inbox domaininbox.Inbox
	err := row.Scan(&inbox.ID, &inbox.UserID, &inbox.MessageID, &inbox.CreatedAt)
	return &inbox, err
}

func (r *inboxRepository) UpdateInbox(ctx context.Context, inbox *domaininbox.Inbox) error {
	_, err := r.pool.Exec(ctx, "UPDATE inboxes SET user_id = $1, message_id = $2, created_at = $3 WHERE id = $4", inbox.UserID, inbox.MessageID, inbox.CreatedAt, inbox.ID)
	return err
}

func (r *inboxRepository) DeleteInbox(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM inboxes WHERE id = $1", id)
	return err
}

func (r *inboxRepository) ListInboxes(ctx context.Context, userID int64, page int, pageSize int) ([]*domaininbox.Inbox, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, user_id, message_id, created_at FROM inboxes WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", userID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inboxes := make([]*domaininbox.Inbox, 0)
	for rows.Next() {
		var inbox domaininbox.Inbox
		err := rows.Scan(&inbox.ID, &inbox.UserID, &inbox.MessageID, &inbox.CreatedAt)
		if err != nil {
			return nil, err
		}
		inboxes = append(inboxes, &inbox)
	}
	return inboxes, nil
}
