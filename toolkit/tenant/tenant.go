package tenant

import (
	"context"
	"fmt"
	"ovra/toolkit/auth"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

func GetTenantId(ctx context.Context, rds *redis.Redis, user *auth.UserInfo) (string, error) {
	//获取redis缓存的key
	key := fmt.Sprintf(TENANT_KEY, user.UserId)
	// 先判断redis中是否存在
	ex, err := rds.ExistsCtx(ctx, key)
	if err != nil {
		return "", err
	}
	if !ex {
		return user.TenantId, nil
	}
	val, err := rds.HgetCtx(ctx, key, "nt")
	if err != nil {
		return "", err
	}
	return val, nil
}

func SetTenantId(ctx context.Context, rds *redis.Redis, userId, tenantId string) error {
	key := fmt.Sprintf(TENANT_KEY, userId)
	ot := auth.GetTenantId(ctx)
	ex, err := rds.ExistsCtx(ctx, key)
	if err != nil {
		return err
	}
	if !ex {
		err = rds.HsetCtx(ctx, key, "ot", ot)
		if err != nil {
			return err
		}
	}
	err = rds.HsetCtx(ctx, key, "nt", tenantId)
	if err != nil {
		return err
	}
	return nil
}
