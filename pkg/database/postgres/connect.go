package postgres

import (
	"context"
	"fmt"

	"github.com/Slintox/user-service/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Connect(ctx context.Context, cfg *config.PostgresConfig) (*pgxpool.Pool, error) {
	pgCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	pgPool, err := pgxpool.ConnectConfig(ctx, pgCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get db connection: %w", err)
	}

	if err = pgPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	return pgPool, nil
}
