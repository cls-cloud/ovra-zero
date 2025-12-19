package main

import (
	"flag"
	"fmt"
	"os"
	"ovra/app/system/internal/config"
	"ovra/app/system/internal/handler"
	"ovra/app/system/internal/middleware"
	"ovra/app/system/internal/svc"
	"ovra/toolkit/helper"
	"ovra/toolkit/middlewares"
	"ovra/toolkit/utils"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile *string

func init() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	defaultConfig := filepath.Join("etc", env, "system.yaml")
	configFile = flag.String("f", defaultConfig, "the config file")
	flag.Parse()
	fmt.Println("Using config file:", *configFile)
}

func main() {
	flag.Parse()
	if err := utils.Init("2024-01-01", 1); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建服务器并传入自定义的 UnauthorizedCallback
	server := rest.MustNewServer(c.RestConf, rest.WithCors("*"))

	// 使用拦截器
	httpx.SetOkHandler(helper.OkHandler)
	httpx.SetErrorHandlerCtx(helper.ErrHandler(c.RestConf.Name))

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	// 注册中间件
	server.Use(middleware.LogMiddleware)
	server.Use(middlewares.ApiMiddleware(c.RestConf.Mode))

	//rpc := zrpc.MustNewServer(c.RpcConf, func(grpcServer *grpc.Server) {
	//	system.RegisterUserRpcServer(grpcServer, userServer.NewUserRpcServer(ctx))
	//	system.RegisterRoleRpcServer(grpcServer, roleServer.NewRoleRpcServer(ctx))
	//
	//	if c.RpcConf.Mode == service.DevMode || c.RpcConf.Mode == service.TestMode {
	//		reflection.Register(grpcServer)
	//	}
	//})
	group := service.NewServiceGroup()
	group.Add(server)
	//group.Add(rpc)
	defer group.Stop()
	fmt.Printf("Starting server at %s:%d...\n", c.RestConf.Host, c.RestConf.Port)
	//fmt.Printf("Starting rpc server at %s...\n", c.RpcConf.ListenOn)
	group.Start()
}
