package auth

import (
	"context"

	"github.com/K1la/warehouse-control/internal/dto"
	"github.com/rs/zerolog"
)

type Handler struct {
	service Service
	log     zerolog.Logger
}

func New(s Service, l zerolog.Logger) *Handler {
	return &Handler{service: s, log: l}
}

type Service interface {
	Login(ctx context.Context, req dto.LoginRequest) (string, error)
	Register(ctx context.Context, req dto.RegisterRequest) (int64, error)
}
