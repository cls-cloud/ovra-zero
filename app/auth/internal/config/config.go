// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	RestConf  rest.RestConf
	Data      DataConfig
	SystemRpc zrpc.RpcClientConf
	JwtAuth   struct {
		AccessSecret         string
		AccessExpire         int64
		MultipleLoginDevices bool
	}
	ApiDecrypt struct {
		Enabled    bool
		HeaderFlag string
		PublicKey  string
		PrivateKey string
	}
	Captcha struct {
		Enabled bool
	}
}

type DataConfig struct {
	Redis redis.RedisConf
}
