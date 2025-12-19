package main

import (
	"flag"
	"fmt"
	"os"
	"ovra/toolkit/helper"
	"ovra/toolkit/utils"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest/httpx"

	"ovra/app/resource/internal/config"
	"ovra/app/resource/internal/handler"
	"ovra/app/resource/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile *string

func init() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	defaultConfig := filepath.Join("etc", env, "resource.yaml")
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

	group := service.NewServiceGroup()
	group.Add(server)
	defer group.Stop()
	fmt.Printf("Starting server at %s:%d...\n", c.RestConf.Host, c.RestConf.Port)
	group.Start()
}
