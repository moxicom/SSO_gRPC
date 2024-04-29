package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/moxicom/SSO_gRPC/internal/app/grpc"
	authgrpc "github.com/moxicom/SSO_gRPC/internal/grpc/auth"
)

type App struct {
	GRPCSrv *grpcapp.App
}

type Mock struct {
	authgrpc.Auth
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// TODO: init storage

	// TODO: init server layer

	grpcApp := grpcapp.New(log, Mock{}, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
