package svc

import (
	"report/internal/config"
	"report/internal/dal"
	"report/internal/dal/query"
	"report/internal/svc/dbs"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Db     *gorm.DB
	Rds    *redis.Redis
	Dal    *dal.Dal
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := dbs.NewDb(c)
	return &ServiceContext{
		Config: c,
		Rds:    redis.MustNewRedis(c.Data.Redis),
		Db:     db,
		Dal:    dal.NewDal(db, query.Use(db), c),
	}
}
