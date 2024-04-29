package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/moxicom/SSO_gRPC/internal/domain/models"
	"github.com/moxicom/SSO_gRPC/storage"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

// Interact with storage
type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int) (bool, error)
}

// Interact with storage
type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// New creates a new instance of the Auth struct.
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    userSaver,
		usrProvider: userProvider,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

// Login authenticates a user and generates a token for further access.
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (token string, err error) {
	const op = "auth.Login"

	log := a.log.With(slog.String("op", op), slog.String("email", email))
	log.Info("logging in")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("invalid credentials", err.Error())
			return "", fmt.Errorf("%s : %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", err.Error())
		return "", fmt.Errorf("%s : %w", op, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid password", err.Error())
		return "", fmt.Errorf("%s : %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Info("app not found", err.Error())
			return "", fmt.Errorf("%s : %w", op, storage.ErrAppNotFound)
		}

		a.log.Error("failed to get app", err.Error())
		return "", fmt.Errorf("%s : %w", op, err.Error())
	}
}

// RegisterNewUser registers a new user with the given email and password.
// if user with given username is already exists, returns error
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (userID int64, err error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(slog.String("op", op), slog.String("email", email))
	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to save user", err)
		return 0, fmt.Errorf("%s : %w", op, err.Error())
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", err)
		return 0, fmt.Errorf("%s : %w", op, err.Error())
	}

	return id, nil
}

// IsAdmin checks if the user with the given ID is an admin.
func (a *Auth) IsAdmin(ctx context.Context, userID int) (bool, error) {
	panic("auth service method IsAdmin not implemented")
}
