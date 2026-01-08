// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"net/http"
	"ovra/app/demo/internal/config"
	"ovra/toolkit/middlewares"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type AuthMiddleware struct {
	c   config.Config
	rds *redis.Redis
}

func NewAuthMiddleware(c config.Config, rds *redis.Redis) *AuthMiddleware {
	return &AuthMiddleware{
		c:   c,
		rds: rds,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return middlewares.ExecHandle(next, m.c.JwtAuth.AccessSecret, m.rds, m.c.JwtAuth.MultipleLoginDevices)
}
