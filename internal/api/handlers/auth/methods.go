package items

import (
	"errors"
	"github.com/K1la/warehouse-control/internal/dto"
	serviceuser "github.com/K1la/warehouse-control/internal/service/user"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"net/http"
)

// POST /api/auth/login
func (h *Handler) Login(c *ginext.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error().Err(err).Msg("failed to bind login request")
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}
	h.log.Info().Interface("req", req).Msg("login request received")

	token, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, serviceuser.ErrInvalidCredentials) {
			h.log.Error().Err(err).Msg("invalid credentials")
			c.JSON(http.StatusUnauthorized, ginext.H{"error": "invalid credentials"})
			return
		}

		if errors.Is(err, serviceuser.ErrUserNotFound) {
			h.log.Error().Err(err).Msg("user not found")
			c.JSON(http.StatusNotFound, ginext.H{"error": "user not found"})
			return
		}

		h.log.Error().Err(err).Msg("failed to login")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "internal server error"})
		return
	}

	h.log.Info().Interface("token", token).Msg("login success")
	c.JSON(http.StatusOK, dto.LoginResponse{Token: token})
}

func (h *Handler) Register(c *ginext.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to bind json")
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid request body"})
		return
	}

	id, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, serviceuser.ErrUserAlreadyExists) {
			zlog.Logger.Error().Err(err).Msg("user already exists")
			c.JSON(http.StatusConflict, ginext.H{"error": "user already exists"})
			return
		}

		zlog.Logger.Error().Err(err).Msg("failed to register user")
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"id": id})
}
