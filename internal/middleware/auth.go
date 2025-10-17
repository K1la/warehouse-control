package middleware

import (
	"errors"
	jwtpkg "github.com/K1la/warehouse-control/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"net/http"
	"strings"
)

var (
	ErrNoToken            = errors.New("missing token")
	ErrInvalidTokenHeader = errors.New("invalid token header")
	ErrInvalidTokenFormat = errors.New("invalid token format")
	ErrExpiredToken       = errors.New("token had expired")
	ErrRoleNotFound       = errors.New("role not found in context")
	ErrInvalidRole        = errors.New("invalid role type")
	ErrAccessDenied       = errors.New("access forbidden")
)

type AuthMW struct {
	jwt *jwtpkg.JWT
}

func NewAuth(j *jwtpkg.JWT) *AuthMW {
	return &AuthMW{jwt: j}
}

func (a *AuthMW) RequireAuth() ginext.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrNoToken)
			return
		}
		zlog.Logger.Info().Str("auth", auth).Msg("get header MW")

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrInvalidTokenFormat)
			return
		}

		claims, err := a.jwt.Parse(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		zlog.Logger.Info().Interface("claims", claims).Msg("get claims")

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func (a *AuthMW) RequireRole(roles ...string) ginext.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *ginext.Context) {
		r, ok := c.Get("role")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrRoleNotFound)
			return
		}
		role, ok := r.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrInvalidRole)
			return
		}

		if _, ok = allowed[role]; !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrAccessDenied)
			return
		}

		c.Next()
	}
}
