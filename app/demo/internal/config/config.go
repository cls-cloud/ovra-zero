// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	RestConf rest.RestConf
	RpcConf  zrpc.RpcServerConf
	Tenant   TenantConfig
	Data     DataConfig
	JwtAuth  struct {
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

type TenantConfig struct {
	Enabled      bool
	IgnoreTables []string
}

type DataConfig struct {
	Database DatabaseConfig
	Redis    redis.RedisConf
	Cache    CacheConfig
}
type DatabaseConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
}
type CacheConfig struct {
	Expire int
}
