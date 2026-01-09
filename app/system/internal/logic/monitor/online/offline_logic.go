package online

import (
	"context"
	"fmt"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"ovra/toolkit/auth"
	"ovra/toolkit/errx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OfflineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOfflineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OfflineLogic {
	return &OfflineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OfflineLogic) Offline(req *types.IdReq) error {
	token := req.Id
	us, err := auth.AnalyseToken(token, l.svcCtx.Config.JwtAuth.AccessSecret)
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
