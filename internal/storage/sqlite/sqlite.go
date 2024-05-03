package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
	"github.com/moxicom/SSO_gRPC/internal/domain/models"
	"github.com/moxicom/SSO_gRPC/internal/storage"
	"github.com/moxicom/SSO_gRPC/internal/utils"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, utils.MakeError(op, err)
	}

	return &Storage{db: db}, nil
}

// Потом еще раз глянуть про контекст надо
// SaveUser saves a user by email and password hash
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?, ?)")
	if err != nil {
		return 0, utils.MakeError(op, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, utils.MakeError(op, storage.ErrUserAlreadyExists)
		}

		return 0, utils.MakeError(op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, utils.MakeError(op, err)
	}

	return id, nil
}

// User returns user by email
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, utils.MakeError(op, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, utils.MakeError(op, storage.ErrUserNotFound)
		}

		return models.User{}, utils.MakeError(op, err)
	}

	return user, nil
}

// IsAdmin returns is user an admin
func (s *Storage) IsAdmin(ctx context.Context, userID uint64) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, utils.MakeError(op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool
	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, utils.MakeError(op, storage.ErrAppNotFound)
		}

		return false, utils.MakeError(op, err)
	}

	return isAdmin, nil
}
