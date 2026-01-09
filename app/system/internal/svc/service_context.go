package svc

import (
	"ovra/app/system/internal/config"
	"ovra/app/system/internal/dal"
	"ovra/app/system/internal/dal/query"
	"ovra/app/system/internal/middleware"
	"ovra/app/system/internal/svc/dbs"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Db     *gorm.DB
	Rds    *redis.Redis
	Query  *query.Query
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
		Query:  query.Use(db),
		Auth:   middleware.NewAuthMiddleware(c, rds).Handle,
		Dal:    dal.NewDal(db, query.Use(db), c),
	}
}
