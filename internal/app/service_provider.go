package app

import (
	"context"
	"github.com/Slintox/user-service/pkg/common/closer"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"

	"github.com/Slintox/user-service/config"
	userImpl "github.com/Slintox/user-service/internal/api/user"
	userRepo "github.com/Slintox/user-service/internal/repository/user"
	userService "github.com/Slintox/user-service/internal/service/user"
	"github.com/Slintox/user-service/pkg/database/postgres"
)

type ServiceProvider interface {
	GetConfig() *config.Config
	GetPostgresClient(ctx context.Context) postgres.Client

	GetUserRepository(ctx context.Context) userRepo.Repository
	GetUserService(ctx context.Context) userService.Service
	GetUserImpl(ctx context.Context) *userImpl.Implementation
}

type serviceProvider struct {
	configPath string

	config   *config.Config
	pgClient postgres.Client

	userRepo    userRepo.Repository
	userService userService.Service

	userImpl *userImpl.Implementation
}

func newServiceProvider(configPath string) *serviceProvider {
	return &serviceProvider{
		configPath: configPath,
	}
}

func (s *serviceProvider) GetConfig() *config.Config {
	if s.config != nil {
		return s.config
	}

	var err error
	s.config, err = config.InitConfig(s.configPath)
	if err != nil {
		log.Fatalf("failed to init config: %s", err.Error())
		return nil
	}

	return s.config
}

func (s *serviceProvider) GetPostgresClient(ctx context.Context) postgres.Client {
	if s.pgClient != nil {
		return s.pgClient
	}

	pgCfg, err := pgxpool.ParseConfig(s.GetConfig().GetPostgresConfig().DSN)
	if err != nil {
		log.Fatalf("failed to get db config: %s", err.Error())
	}

	client, err := postgres.NewClient(ctx, pgCfg)
	if err != nil {
		log.Fatalf("failed to get postgres client: %s", err.Error())
	}

	err = client.Postgres().Ping(ctx)
	if err != nil {
		log.Fatalf("ping error: %s", err.Error())
	}
	closer.Add(client.Close)

	s.pgClient = client
	return s.pgClient
}

func (s *serviceProvider) GetUserRepository(ctx context.Context) userRepo.Repository {
	if s.userRepo != nil {
		return s.userRepo
	}

	s.userRepo = userRepo.NewRepository(s.GetPostgresClient(ctx))
	return s.userRepo
}

func (s *serviceProvider) GetUserService(ctx context.Context) userService.Service {
	if s.userService != nil {
		return s.userService
	}

	s.userService = userService.NewService(s.GetUserRepository(ctx))
	return s.userService
}

func (s *serviceProvider) GetUserImpl(ctx context.Context) *userImpl.Implementation {
	if s.userImpl != nil {
		return s.userImpl
	}

	s.userImpl = userImpl.NewImplementation(s.GetUserService(ctx))
	return s.userImpl
}
