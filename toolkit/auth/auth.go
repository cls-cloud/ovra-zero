package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	// 上下文字段
	UserIDKey       = "userId"
	TenantIDKey     = "tenantId"
	ClientIDKey     = "clientId"
	DeptNameKey     = "deptName"
	DataScopeKey    = "dataScope"
	CurrentDeptKey  = "currentDept"
	BellowDeptKey   = "bellowDept"
	CustomerDeptKey = "customerDept"

	// Redis key 模板
	TokenKey    = "token:%s:%s"    // clientId + userId
	TokenKeyMd5 = "token:%s:%s:%s" // clientId + userId + md5

	// Redis Hash 字段
	FieldToken         = "token"
	FieldActiveTimeout = "activeTimeout"
	FieldCurrentTime   = "currentTime"
	FieldExpireTime    = "expireTime"
	FieldLoginInfoId   = "loginInfoId"
)

type Auth struct {
	rds  *redis.Redis
	user *UserInfo
}

func NewAuth(rds *redis.Redis, user *UserInfo) *Auth {
	return &Auth{rds: rds, user: user}
}

// SetToken 保存登录 token 信息到 Redis（含过期与滑动窗口时间）
func (a *Auth) SetToken(ctx context.Context, key, token string, activeTimeout, ttl int64, loginInfoId string) error {
	now := time.Now().Unix()

	fields := map[string]string{
		FieldToken:         token,
		FieldActiveTimeout: strconv.FormatInt(activeTimeout, 10),
		FieldCurrentTime:   strconv.FormatInt(now, 10),
		FieldLoginInfoId:   loginInfoId,
	}

	for f, v := range fields {
		if err := a.rds.HsetCtx(ctx, key, f, v); err != nil {
			return fmt.Errorf("set %s failed: %v", f, err)
		}
	}
	if err := a.rds.ExpireCtx(ctx, key, int(time.Duration(ttl)*time.Second)); err != nil {
		return fmt.Errorf("set expire failed: %v", err)
	}
	return nil
}

// CheckToken 检查 token 是否活跃，并刷新 activeTimeout
func (a *Auth) CheckToken(ctx context.Context, key, tokenStr string) (bool, error) {
	exists, err := a.rds.ExistsCtx(ctx, key)
	if err != nil {
		return true, fmt.Errorf("check key exist failed: %v", err)
	}
	if !exists {
		return true, fmt.Errorf("token key does not exist")
	}

	tkStr, err := a.rds.HgetCtx(ctx, key, FieldToken)
	if err != nil {
		return true, fmt.Errorf("get token failed: %v", err)
	}
	if tkStr != tokenStr {
		return true, fmt.Errorf("invalid token")
	}
	curStr, err := a.rds.HgetCtx(ctx, key, FieldCurrentTime)
	if err != nil {
		return true, fmt.Errorf("get currentTime failed: %v", err)
	}
	actStr, err := a.rds.HgetCtx(ctx, key, FieldActiveTimeout)
	if err != nil {
		return true, fmt.Errorf("get activeTimeout failed: %v", err)
	}

	curInt, _ := strconv.ParseInt(curStr, 10, 64)
	actInt, _ := strconv.ParseInt(actStr, 10, 64)
	now := time.Now().Unix()

	if curInt > 0 && actInt > 0 && now > curInt+actInt {
		return true, nil // token 过期
	}

	if err := a.rds.HsetCtx(ctx, key, FieldCurrentTime, strconv.FormatInt(now, 10)); err != nil {
		return false, fmt.Errorf("refresh currentTime failed: %v", err)
	}

	return false, nil
}
