package app

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Slintox/user-service/config"
	"github.com/Slintox/user-service/internal/api/user"
	uRepo "github.com/Slintox/user-service/internal/repository/user"
	uService "github.com/Slintox/user-service/internal/service/user"
	"github.com/Slintox/user-service/pkg/database/postgres"
	userV1 "github.com/Slintox/user-service/pkg/user_v1"
)

func Run(configPath string) {
	ctx := context.Background()

	cfg, err := config.InitConfig(configPath)
	if err != nil {
		log.Fatalf("failed to get config: %s", err.Error())
	}

	list, err := net.Listen("tcp", cfg.GRPC.Port)
	if err != nil {
		log.Fatalf("failed to get listener: %s", err.Error())
	}

	s := grpc.NewServer()
	reflection.Register(s)

	pgPool, err := postgres.Connect(ctx, cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to get postgres connect: %s", err.Error())
	}

	var userRepo uRepo.Repository
	var userService uService.Service

	userRepo = uRepo.NewRepository(pgPool)
	userService = uService.NewService(userRepo)
	userV1.RegisterUserV1Server(s, user.NewImplementation(userService))

	if err = s.Serve(list); err != nil {
		log.Fatalf("failed to serve: %s", err.Error())
	}
}
