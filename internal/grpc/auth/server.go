//
// Contains auth gGRPC handlers here
//

package auth

import (
	"context"

	ssov1 "github.com/moxicom/SSO_gRPC_PROTOS/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int, err error)
	IsAdmin(ctx context.Context, userID int) (bool, error)
}

type ServerAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

const (
	emptyValue = 0
)

func Register(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &ServerAPI{auth: auth})
}

func (s *ServerAPI) Login(
	ctx context.Context, r *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLogin(r); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, r.GetEmail(), r.GetPasword(), int(r.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) Register(
	ctx context.Context, r *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(r); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, r.GetEmail(), r.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: int64(userID),
	}, nil
}

func (s *ServerAPI) IsAdmin(
	ctx context.Context, r *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	panic("method IsAdmin not implemented")
}

func validateLogin(r *ssov1.LoginRequest) error {
	if r.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if r.GetPasword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if r.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "appId is required")
	}

	return nil
}

func validateRegister(r *ssov1.RegisterRequest) error {
	if r.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if r.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}
