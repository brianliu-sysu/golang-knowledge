package postgres

import (
	"context"

	domainmessage "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/domain/message"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/respository/message"
	"github.com/jackc/pgx/v5/pgxpool"
)

type messageRepository struct {
	pool *pgxpool.Pool
}

func NewMessageRepository(pool *pgxpool.Pool) message.MessageRepository {
	return &messageRepository{pool: pool}
}

func (r *messageRepository) CreateMessage(ctx context.Context, message *domainmessage.Message) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO messages (id, from_user, to_user, group_id, content, created_at) VALUES ($1, $2, $3, $4, $5, $6)", message.ID, message.FromUser, message.ToUser, message.GroupID, message.Content, message.CreatedAt)
	return err
}

func (r *messageRepository) GetMessage(ctx context.Context, id int64) (*domainmessage.Message, error) {
	row := r.pool.QueryRow(ctx, "SELECT id, from_user, to_user, group_id, content, created_at FROM messages WHERE id = $1", id)
	var message domainmessage.Message
	err := row.Scan(&message.ID, &message.FromUser, &message.ToUser, &message.GroupID, &message.Content, &message.CreatedAt)
	return &message, err
}

func (r *messageRepository) UpdateMessage(ctx context.Context, message *domainmessage.Message) error {
	_, err := r.pool.Exec(ctx, "UPDATE messages SET from_user = $1, to_user = $2, group_id = $3, content = $4, created_at = $5 WHERE id = $6", message.FromUser, message.ToUser, message.GroupID, message.Content, message.CreatedAt, message.ID)
	return err
}

func (r *messageRepository) DeleteMessage(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM messages WHERE id = $1", id)
	return err
}

func (r *messageRepository) ListMessages(ctx context.Context, userID int64, page int, pageSize int) ([]*domainmessage.Message, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, from_user, to_user, group_id, content, created_at FROM messages WHERE from_user = $1 OR to_user = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", userID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*domainmessage.Message, 0)
	for rows.Next() {
		var message domainmessage.Message
		err := rows.Scan(&message.ID, &message.FromUser, &message.ToUser, &message.GroupID, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}
	return messages, nil
}

func (r *messageRepository) ListMessagesByGroupID(ctx context.Context, groupID int64, page int, pageSize int) ([]*domainmessage.Message, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, from_user, to_user, group_id, content, created_at FROM messages WHERE group_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", groupID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*domainmessage.Message, 0)
	for rows.Next() {
		var message domainmessage.Message
		err := rows.Scan(&message.ID, &message.FromUser, &message.ToUser, &message.GroupID, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}
	return messages, nil
}
