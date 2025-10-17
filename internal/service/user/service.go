package user

import (
	"context"

	"github.com/K1la/warehouse-control/internal/model"
	jwtpkg "github.com/K1la/warehouse-control/pkg/jwt"

	"github.com/rs/zerolog"
)

type Service struct {
	db  Repo
	j   *jwtpkg.JWT
	log zerolog.Logger
}

func New(d Repo, jwt *jwtpkg.JWT, l zerolog.Logger) *Service {
	return &Service{db: d, j: jwt, log: l}
}

type Repo interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	CheckUserExistByUsername(ctx context.Context, username string) (bool, error)
}
