// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"ovra/app/demo/internal/config"
	"ovra/app/demo/internal/dal"
	"ovra/app/demo/internal/dal/query"
	"ovra/app/demo/internal/middleware"
	"ovra/app/demo/internal/svc/dbs"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Db     *gorm.DB
	Rds    *redis.Redis
	Auth   rest.Middleware
	Dal    *dal.Dal
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := dbs.NewDb(c)
	rds := redis.MustNewRedis(c.Data.Redis)
	return &ServiceContext{
		Config: c,
		Rds:    rds,
		Db:     db,
		Auth:   middleware.NewAuthMiddleware(c, rds).Handle,
		Dal:    dal.NewDal(db, query.Use(db), c),
	}
}
