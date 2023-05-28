package app

import (
	"context"
	"github.com/Slintox/user-service/pkg/grpc/interceptor/validate"
	"github.com/Slintox/user-service/pkg/service/user_v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	statikFs "github.com/rakyll/statik/fs"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Slintox/user-service/pkg/common/closer"
	_ "github.com/Slintox/user-service/statik"
)

type App struct {
	configPath string

	serviceProvider ServiceProvider

	grpcServer *grpc.Server

	// HTTP server for the gRPC gateway
	httpGatewayServer *http.Server
	httpSwaggerServer *http.Server
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

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := a.runGRPCServer(); err != nil {
			log.Fatalf("failed to run GRPC server: %s", err.Error())
		}

	}()

	wg.Add(1)
	go func() {
		if err := a.runHTTPServer(); err != nil {
			log.Fatalf("failed to run HTTP gateway server: %s", err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		if err := a.runSwaggerServer(); err != nil {
			log.Fatalf("failed to run Swagger server: %s", err.Error())
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
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
	a.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(
			validate.UnaryInterceptor,
		))
	reflection.Register(a.grpcServer)

	user_v1.RegisterUserV1Server(a.grpcServer, a.serviceProvider.GetUserImpl(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := user_v1.RegisterUserV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GetConfig().GetGRPCConfig().Port, dialOpts)
	if err != nil {
		return err
	}

	a.httpGatewayServer = &http.Server{
		Addr:    a.serviceProvider.GetConfig().GetHTTPConfig().Port,
		Handler: mux,
	}

	return nil
}

func (a *App) initSwaggerServer(ctx context.Context) error {
	fs, err := statikFs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(fs)))
	mux.HandleFunc("/api.swagger.json", a.serveSwaggerFile("/api.swagger.json"))

	a.httpSwaggerServer = &http.Server{
		Addr:    a.serviceProvider.GetConfig().GetSwaggerConfig().Port,
		Handler: mux,
	}

	return nil
}

func (a *App) runGRPCServer() error {
	list, err := net.Listen("tcp", a.serviceProvider.GetConfig().GetGRPCConfig().Port)
	if err != nil {
		log.Fatalf("failed to get listener: %s", err.Error())
	}

	log.Printf("GRPC server is running on port %s", a.serviceProvider.GetConfig().GetGRPCConfig().Port)

	if err = a.grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve: %s", err.Error())
	}

	return nil
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP gateway server is running on port %s", a.serviceProvider.GetConfig().GetHTTPConfig().Port)

	if err := a.httpGatewayServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	log.Printf("Swagger server is running on port %s", a.serviceProvider.GetConfig().GetSwaggerConfig().Port)

	if err := a.httpSwaggerServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (a *App) serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("serving swagger file: %s", path)

		fs, err := statikFs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		file, err := fs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
