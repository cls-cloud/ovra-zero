package svc

import (
	"monitor/client/logininforpc"
	"monitor/client/operlogrpc"
	"monitor/pb/monitor"
	"system/internal/config"
	"system/internal/dal"
	q "system/internal/dal/query"
	"system/internal/dao/query"
	"system/internal/middleware"
	"system/internal/svc/dbs"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	Db           *gorm.DB
	Rds          *redis.Redis
	Query        *query.Query
	Auth         rest.Middleware
	LoginInfoRpc monitor.LoginInfoRpcClient
	OperLogRpc   monitor.OperLogRpcClient
	Dal          *dal.Dal
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := dbs.NewDb(c)
	rds := redis.MustNewRedis(c.Data.Redis)
	return &ServiceContext{
		Config:       c,
		Rds:          rds,
		Db:           db,
		Query:        query.Use(db),
		Auth:         middleware.NewAuthMiddleware(c, rds).Handle,
		LoginInfoRpc: logininforpc.NewLoginInfoRpc(zrpc.MustNewClient(c.MonitorRpc)),
		OperLogRpc:   operlogrpc.NewOperLogRpc(zrpc.MustNewClient(c.MonitorRpc)),
		Dal:          dal.NewDal(db, q.Use(db), c),
	}
}
