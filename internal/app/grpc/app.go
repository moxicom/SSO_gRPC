package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/moxicom/SSO_gRPC/internal/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	authService authgrpc.Auth,
	port int,
) *App {

	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// Run a gRPC server and panics if any error occurs
func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

// Run a gRPC server
func (a *App) run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op))
	log.Info("Starting gRPC server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	reflection.Register(a.gRPCServer)

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop a gRPC server
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("Stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
