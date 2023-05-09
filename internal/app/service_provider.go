package app

import (
	"context"
	"log"

	"github.com/Slintox/user-service/config"
	userImpl "github.com/Slintox/user-service/internal/api/user"
	userRepo "github.com/Slintox/user-service/internal/repository/user"
	userService "github.com/Slintox/user-service/internal/service/user"
	"github.com/Slintox/user-service/pkg/database/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ServiceProvider interface {
}

type serviceProvider struct {
	configPath string

	config *config.Config
	pgConn *pgxpool.Pool

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

func (s *serviceProvider) GetPostgresClient(ctx context.Context) *pgxpool.Pool {
	if s.pgConn != nil {
		return s.pgConn
	}

	conn, err := postgres.Connect(ctx, s.GetConfig().GetPostgresConfig())
	if err != nil {
		log.Fatalf("falied to connect to postgres db: %s", err.Error())
		return nil
	}

	s.pgConn = conn
	return s.pgConn
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
