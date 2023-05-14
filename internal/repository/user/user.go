package user

import (
	"context"
	"errors"
	"github.com/Slintox/user-service/pkg/database/postgres"
	"github.com/jackc/pgx/v4"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/Slintox/user-service/config"
	"github.com/Slintox/user-service/internal/model"
	repo "github.com/Slintox/user-service/internal/repository"
)

const tableName = `"user"`

type Repository interface {
	Add(ctx context.Context, user *model.CreateUser) error
	Get(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, username string, updateData *model.UpdateUser) error
	Delete(ctx context.Context, username string) error
	IsUsernameAvailable(ctx context.Context, username string) (bool, error)
}

type repository struct {
	client postgres.Client
}

func NewRepository(client postgres.Client) Repository {
	return &repository{
		client: client,
	}
}

func (r *repository) Add(ctx context.Context, user *model.CreateUser) error {
	var roleId int

	roleQ := postgres.Query{
		Name:     "userRole.Get",
		QueryRaw: "select id from user_role where id = $1",
	}

	row := r.client.Postgres().QueryRow(ctx, roleQ, user.Role)
	if err := row.Scan(&roleId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.ErrRecordNotFound
		}
		return err
	}

	builder := sq.Insert(tableName).
		Columns("username", "email", "password", "role").
		Values(user.Username, user.Email, user.Password, user.Role).
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	if config.PostgresDev {
		log.Printf("user.Update: query: '%s' values: '%+v'\n", query, v)
	}

	q := postgres.Query{
		Name:     "user.Add",
		QueryRaw: query,
	}

	_, err = r.client.Postgres().Exec(ctx, q, v...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Get(ctx context.Context, username string) (*model.User, error) {
	builder := sq.Select("username", "email", "password", "role", "created_at", "updated_at").
		From(tableName).
		Where(sq.Eq{"username": username}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	if config.PostgresDev {
		log.Printf("user.Get: query: '%s' values: '%+v'\n", query, v)
	}

	q := postgres.Query{
		Name:     "user.Get",
		QueryRaw: query,
	}

	var user model.User
	if err = r.client.Postgres().Get(ctx, &user, q, v...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) Update(ctx context.Context, username string, updateData *model.UpdateUser) error {
	updateQuery := sq.Update(tableName).
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar)

	if updateData.Username != nil {
		updateQuery = updateQuery.Set("username", updateData.Username)
	}
	if updateData.Password != nil {
		updateQuery = updateQuery.Set("password", updateData.Password)
	}
	if updateData.Email != nil {
		updateQuery = updateQuery.Set("email", updateData.Email)
	}
	if updateData.Role != nil {
		updateQuery = updateQuery.Set("role", updateData.Role)
	}

	updateQuery = updateQuery.Set("updated_at", "now()")

	query, v, err := updateQuery.ToSql()
	if err != nil {
		return err
	}

	if config.PostgresDev {
		log.Printf("user.Update: query: '%s' values: '%+v'\n", query, v)
	}

	q := postgres.Query{
		Name:     "user.Update",
		QueryRaw: query,
	}

	pg, err := r.client.Postgres().Exec(ctx, q, v...)
	if err != nil {
		return err
	}

	if pg.RowsAffected() == 0 {
		return repo.ErrRecordNotFound
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, username string) error {
	builder := sq.Delete(tableName).
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	if config.PostgresDev {
		log.Printf("user.Delete: query: '%s' values: '%+v'\n", query, v)
	}

	q := postgres.Query{
		Name:     "user.Delete",
		QueryRaw: query,
	}

	_, err = r.client.Postgres().Exec(ctx, q, v...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) IsUsernameAvailable(ctx context.Context, username string) (bool, error) {
	builder := sq.Select("count(*)").
		From(tableName).
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return false, err
	}

	if config.PostgresDev {
		log.Printf("user.IsUsernameAvailable: query: '%s' values: '%+v'\n", query, v)
	}

	q := postgres.Query{
		Name:     "user.IsUsernameAvailable",
		QueryRaw: query,
	}

	var count int
	row := r.client.Postgres().QueryRow(ctx, q, v...)
	if err = row.Scan(&count); err != nil {
		log.Printf("user.IsUnameAvl: %s", err.Error())
		return false, err
	}

	return count == 0, nil
}
