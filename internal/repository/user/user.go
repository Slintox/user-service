package user

import (
	"context"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/Slintox/user-service/internal/model"
	repo "github.com/Slintox/user-service/internal/repository"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) Add(ctx context.Context, user *model.CreateUser) error {
	var roleId int

	row := r.pool.QueryRow(ctx, "select id from user_role where id = $1", user.Role)
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

	log.Printf("user.Update: query: '%s' values: '%+v'\n", query, v)

	_, err = r.pool.Exec(ctx, query, v...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Get(ctx context.Context, username string) (*model.User, error) {
	builder := sq.Select("username", "email", "password", "role", "created_at", "updated_at").
		From(tableName).
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	log.Printf("user.Get: query: '%s' values: '%+v'\n", query, v)

	rows := r.pool.QueryRow(ctx, query, v...)
	if err != nil {
		return nil, err
	}

	var user model.User
	if err = rows.Scan(&user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
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

	log.Printf("user.Update: query: '%s' values: '%+v'\n", query, v)

	pg, err := r.pool.Exec(ctx, query, v...)
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
		Where("username = $1", username).
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	log.Printf("user.Delete: query: '%s' values: '%+v'\n", query, v)

	_, err = r.pool.Exec(ctx, query, v...)
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

	log.Printf("user.IsUsernameAvailable: query: '%s' values: '%+v'\n", query, v)

	var count int
	row := r.pool.QueryRow(ctx, query, v...)
	if err = row.Scan(&count); err != nil {
		log.Printf("user.IsUnameAvl: %s", err.Error())
		return false, err
	}

	return count == 0, nil
}
