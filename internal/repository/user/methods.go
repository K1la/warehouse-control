package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/K1la/warehouse-control/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

func (r *Postgres) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	query := `
	INSERT INTO users (username, password_hash, role)
	VALUES ($1, $2, $3)
	RETURNING id;
	`

	err := r.db.QueryRowContext(ctx, query, user.Username, user.PasswordHash, user.Role).Scan(
		&user.ID,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return -1, fmt.Errorf("failed create user: %w", err)
	}
	return user.ID, nil
}

func (r *Postgres) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
	SELECT id, username, password_hash, role
	FROM users
	WHERE username = $1
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

func (r *Postgres) CheckUserExistByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	return exists, nil
}
