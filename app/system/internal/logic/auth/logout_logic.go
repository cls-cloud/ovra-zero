package auth

import (
	"context"
	"fmt"
	"net/http"
	"ovra/toolkit/auth"
	"ovra/toolkit/errx"
	"strings"

	"ovra/app/system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *LogoutLogic) Logout() error {
	authorization := l.r.Header.Get("Authorization")
	if authorization == "" {
		return errx.BizErr("Unauthorized: missing token")
	}
	tokenString := strings.TrimPrefix(authorization, "Bearer ")
	us, err := auth.AnalyseToken(tokenString, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		return errx.BizErr("系统异常")
	}
	key := ""
	if l.svcCtx.Config.JwtAuth.MultipleLoginDevices {
		key = fmt.Sprintf(auth.TokenKeyMd5, us.ClientId, us.UserId, us.UsMd5)
	} else {
		key = fmt.Sprintf(auth.TokenKey, us.ClientId, us.UserId)
	}
	_, err = l.svcCtx.Rds.DelCtx(l.ctx, key)
	if err != nil {
		return err
	}
	return nil
}
