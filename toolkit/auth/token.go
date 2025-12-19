package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserInfo
	jwt.RegisteredClaims
}

// GenerateToken 生成 token
func GenerateToken(user UserInfo, secretKey string, expireInSeconds int64) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expireInSeconds) * time.Second)
	UserClaim := &UserClaims{
		UserInfo: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken 解析 token
func AnalyseToken(tokenString, secretKey string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return userClaim, nil
}

// GetUUID 生成唯一码
func GetUUID() string {
	return uuid.New().String()
}

func GetData(ctx context.Context, key string) string {
	val := ctx.Value(key)
	str, ok := val.(string)
	if str == "" || !ok {
		fmt.Printf("============> GetUserId Error")
		return ""
	}
	return str

}
func GetUserId(ctx context.Context) string {
	val := ctx.Value(UserIDKey)
	userIdStr, ok := val.(string)
	if userIdStr == "" || !ok {
		fmt.Printf("============> GetUserId Error")
		return ""
	}
	return userIdStr
}

func GetTenantId(ctx context.Context) string {
	val := ctx.Value(TenantIDKey)
	tenantIdStr, ok := val.(string)
	if tenantIdStr == "" || !ok {
		fmt.Printf("============> GetTenantId Error")
		return ""
	}
	return tenantIdStr
}

func GetUserIdInt(ctx context.Context) int64 {
	userIdStr := GetUserId(ctx)
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	return userId
}
