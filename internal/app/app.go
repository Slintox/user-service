package app

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Slintox/user-service/pkg/common/closer"
	userV1 "github.com/Slintox/user-service/pkg/user_v1"
)

type App struct {
	configPath string

	serviceProvider ServiceProvider

	grpcServer *grpc.Server
}

func NewApp(ctx context.Context, configPath string) (*App, error) {
	app := &App{
		configPath: configPath,
	}

	if err := app.initDeps(ctx); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	if err := a.runGRPCServer(); err != nil {
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider(a.configPath)
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer()
	reflection.Register(a.grpcServer)

	userV1.RegisterUserV1Server(a.grpcServer, a.serviceProvider.GetUserImpl(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	list, err := net.Listen("tcp", a.serviceProvider.GetConfig().GetGRPCConfig().Port)
	if err != nil {
		log.Fatalf("failed to get listener: %s", err.Error())
	}

	if err = a.grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve: %s", err.Error())
	}

	return nil
}
