package user_role

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/Slintox/user-service/internal/model"
	repo "github.com/Slintox/user-service/internal/repository"
	"github.com/Slintox/user-service/pkg/database/postgres"
	"github.com/jackc/pgx/v4"
)

var tableName = "user_role"

type Repository interface {
	// Add добавляет новую роль пользователя.
	Add(ctx context.Context, roleName string) error

	// Get возвращает роль пользователя по её идентификатору.
	Get(ctx context.Context, roleID int) (*model.UserRole, error)

	// Delete удаляет роль пользователя по её идентификатору.
	Delete(ctx context.Context, roleID int) error

	// IsRoleExist проверяет, существует ли роль пользователя по её идентификатору.
	IsRoleExist(ctx context.Context, roleID int) (bool, error)
}

type repository struct {
	client postgres.Client
}

func NewRepository(client postgres.Client) Repository {
	return &repository{
		client: client,
	}
}

// Add adds a new user.
func (r *repository) Add(ctx context.Context, roleName string) error {
	builder := sq.Insert(tableName).
		Columns("name").
		Values(roleName).
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := postgres.Query{
		Name:     "userRole.Add",
		QueryRaw: query,
	}

	_, err = r.client.Postgres().Exec(ctx, q, v...)
	if err != nil {
		return err
	}

	return nil
}

// Get gets a user.
func (r *repository) Get(ctx context.Context, roleID int) (*model.UserRole, error) {
	builder := sq.Select("id", "name").
		From(tableName).
		Where(sq.Eq{"id": roleID}).
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := postgres.Query{
		Name:     "userRole.Get",
		QueryRaw: query,
	}

	userRole := &model.UserRole{}
	if err = r.client.Postgres().Get(ctx, userRole, q, v...); errors.Is(err, pgx.ErrNoRows) {
		return nil, repo.ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return nil, err
}

// Delete deletes a user.
func (r *repository) Delete(ctx context.Context, roleID int) error {
	builder := sq.Delete(tableName).
		Where(sq.Eq{"id": roleID}).
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := postgres.Query{
		Name:     "userRole.Delete",
		QueryRaw: query,
	}

	_, err = r.client.Postgres().Exec(ctx, q, v...)
	if err != nil {
		return err
	}

	return nil
}

// IsRoleExist checks if a role exists.
func (r *repository) IsRoleExist(ctx context.Context, roleID int) (bool, error) {
	builder := sq.Select("id").
		From(tableName).
		Where(sq.Eq{"id": roleID}).
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return false, err
	}

	q := postgres.Query{
		Name:     "userRole.IsRoleExist",
		QueryRaw: query,
	}

	userRole := &model.UserRole{}
	if err = r.client.Postgres().Get(ctx, userRole, q, v...); errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
