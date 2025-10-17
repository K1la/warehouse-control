package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/K1la/warehouse-control/internal/dto"
	"github.com/K1la/warehouse-control/internal/model"
	repouser "github.com/K1la/warehouse-control/internal/repository/user"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

func (s *Service) Register(ctx context.Context, req dto.RegisterRequest) (int64, error) {
	exists, err := s.db.CheckUserExistByUsername(ctx, req.Username)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, ErrUserAlreadyExists
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	user := model.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Role:         req.Role,
	}

	id, err := s.db.CreateUser(ctx, &user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	user, err := s.db.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, repouser.ErrUserNotFound) {
			return "", ErrUserNotFound
		}
		s.log.Error().Err(err).Str("username", req.Username).Msg("failed to get user by username")
		return "", ErrInvalidCredentials
	}

	if err = verifyPassword(req.Password, user.PasswordHash); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.j.Generate(user.ID, user.Username, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// hashPassword generates a bcrypt hash for the given password.
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// verifyPassword checks if the given password matches the stored hash.
func verifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
