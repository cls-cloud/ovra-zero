// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"ovra/app/auth/internal/config"
	"ovra/app/auth/internal/middleware"
	"ovra/app/system/client/sysrpc"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	Rds       *redis.Redis
	Auth      rest.Middleware
	SysClient sysrpc.SysRpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds := redis.MustNewRedis(c.Data.Redis)
	return &ServiceContext{
		Config:    c,
		Rds:       rds,
		Auth:      middleware.NewAuthMiddleware(c, rds).Handle,
		SysClient: sysrpc.NewSysRpc(zrpc.MustNewClient(c.SystemRpc)),
	}
}
