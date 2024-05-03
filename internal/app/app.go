package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/moxicom/SSO_gRPC/internal/app/grpc"
	authgrpc "github.com/moxicom/SSO_gRPC/internal/grpc/auth"
	"github.com/moxicom/SSO_gRPC/internal/services/auth"
	"github.com/moxicom/SSO_gRPC/internal/storage/sqlite"
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
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
