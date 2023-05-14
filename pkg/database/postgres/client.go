package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Client interface {
	Postgres() Postgres
	Close() error
}

var _ Client = &client{}

type client struct {
	pg Postgres
}

func NewClient(ctx context.Context, pgCfg *pgxpool.Config) (*client, error) {
	dbc, err := pgxpool.ConnectConfig(ctx, pgCfg)
	if err != nil {
		log.Fatalf("failed to get db connection: %s", err.Error())
	}

	return &client{
		pg: &postgres{pgxPool: dbc},
	}, nil
}

func (c *client) Postgres() Postgres {
	return c.pg
}

func (c *client) Close() error {
	if &c.pg != nil {
		return c.pg.Close()
	}

	return nil
}
