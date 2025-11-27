package svc

import (
	"monitor/internal/config"
	"monitor/internal/dal"
	"monitor/internal/dal/query"
	"monitor/internal/middleware"
	"monitor/internal/svc/dbs"

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
