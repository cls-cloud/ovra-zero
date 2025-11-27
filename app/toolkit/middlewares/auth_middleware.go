package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"toolkit/auth"
	"toolkit/ip"
	"toolkit/tenant"
	"toolkit/utils"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

func ExecHandle(next http.HandlerFunc, accessSecret string, rds *redis.Redis, multipleLoginDevices bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authorization, "Bearer ")
		uc, err := auth.AnalyseToken(tokenString, accessSecret)
		if err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}
		authInstance := auth.NewAuth(rds, &uc.UserInfo)
		key := ""
		if multipleLoginDevices {
			ipStr, ua := ip.GetIPUa(r)
			name, version := ua.Browser()
			authMd5 := utils.AuthMd5(ipStr, name, version, ua.OS())
			key = fmt.Sprintf(auth.TokenKeyMd5, uc.ClientId, uc.UserId, authMd5)
		} else {
			key = fmt.Sprintf(auth.TokenKey, uc.ClientId, uc.UserId)
		}
		expired, err := authInstance.CheckToken(r.Context(), key, tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}
		if expired {
			http.Error(w, "Unauthorized: token expired (idle timeout)", http.StatusUnauthorized)
			return
		}
		tenantId, err := tenant.GetTenantId(r.Context(), rds, &uc.UserInfo)
		if err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}
		r.Header.Set(auth.UserIDKey, uc.UserId)
		r.Header.Set(auth.TenantIDKey, tenantId)
		r.Header.Set(auth.ClientIDKey, uc.ClientId)
		ctx := context.WithValue(r.Context(), auth.UserIDKey, uc.UserId)
		ctx = context.WithValue(ctx, auth.TenantIDKey, tenantId)
		ctx = context.WithValue(ctx, auth.ClientIDKey, uc.ClientId)
		next(w, r.WithContext(ctx))
	}
}
