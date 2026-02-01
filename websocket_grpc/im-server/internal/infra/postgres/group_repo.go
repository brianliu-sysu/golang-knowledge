package postgres

import (
	"context"

	domaingroup "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/domain/group"
	"github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/respository/group"
	"github.com/jackc/pgx/v5/pgxpool"
)

type groupRepository struct {
	pool *pgxpool.Pool
}

func NewGroupRepository(pool *pgxpool.Pool) group.GroupRepository {
	return &groupRepository{pool: pool}
}

func (r *groupRepository) CreateGroup(ctx context.Context, group *domaingroup.Group) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO groups (id, uuid, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", group.ID, group.UUID, group.Name, group.CreatedAt, group.UpdatedAt)
	return err
}

func (r *groupRepository) GetGroup(ctx context.Context, id int64) (*domaingroup.Group, error) {
	row := r.pool.QueryRow(ctx, "SELECT id, uuid, name, created_at, updated_at FROM groups WHERE id = $1", id)
	var group domaingroup.Group
	err := row.Scan(&group.ID, &group.UUID, &group.Name, &group.CreatedAt, &group.UpdatedAt)
	return &group, err
}

func (r *groupRepository) UpdateGroup(ctx context.Context, group *domaingroup.Group) error {
	_, err := r.pool.Exec(ctx, "UPDATE groups SET uuid = $1, name = $2, created_at = $3, updated_at = $4 WHERE id = $5", group.UUID, group.Name, group.CreatedAt, group.UpdatedAt, group.ID)
	return err
}

func (r *groupRepository) DeleteGroup(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM groups WHERE id = $1", id)
	return err
}

func (r *groupRepository) ListGroups(ctx context.Context, page int, pageSize int) ([]*domaingroup.Group, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, uuid, name, created_at, updated_at FROM groups ORDER BY created_at DESC LIMIT $1 OFFSET $2", pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]*domaingroup.Group, 0)
	for rows.Next() {
		var group domaingroup.Group
		err := rows.Scan(&group.ID, &group.UUID, &group.Name, &group.CreatedAt, &group.UpdatedAt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}
	return groups, nil
}

func (r *groupRepository) ListGroupMembers(ctx context.Context, groupID int64, page int, pageSize int) ([]*domaingroup.GroupMember, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, user_id, group_id, created_at, updated_at FROM group_members WHERE group_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", groupID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groupMembers := make([]*domaingroup.GroupMember, 0)
	for rows.Next() {
		var groupMember domaingroup.GroupMember
		err := rows.Scan(&groupMember.ID, &groupMember.UserID, &groupMember.GroupID, &groupMember.CreatedAt, &groupMember.UpdatedAt)
		if err != nil {
			return nil, err
		}
		groupMembers = append(groupMembers, &groupMember)
	}
	return groupMembers, nil
}
