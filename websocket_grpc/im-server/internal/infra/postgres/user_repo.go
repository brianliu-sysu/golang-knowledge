package postgres

import (
	"context"

	domainuser "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/domain/user"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/respository/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) user.UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domainuser.User) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO users (id, uuid, email, phone, name, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", user.ID, user.UUID, user.Email, user.Phone, user.Name, user.Password, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) GetUser(ctx context.Context, id int64) (*domainuser.User, error) {
	row := r.pool.QueryRow(ctx, "SELECT id, uuid, email, phone, name, password, created_at, updated_at FROM users WHERE id = $1", id)
	var user domainuser.User
	err := row.Scan(&user.ID, &user.UUID, &user.Email, &user.Phone, &user.Name, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	return &user, err
}

func (r *userRepository) UpdateUser(ctx context.Context, user *domainuser.User) error {
	_, err := r.pool.Exec(ctx, "UPDATE users SET uuid = $1, email = $2, phone = $3, name = $4, password = $5, created_at = $6, updated_at = $7 WHERE id = $8", user.UUID, user.Email, user.Phone, user.Name, user.Password, user.CreatedAt, user.UpdatedAt, user.ID)
	return err
}

func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}
