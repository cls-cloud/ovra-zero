package _config

import (
	"context"
	"toolkit/errx"

	"system/internal/svc"
	"system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigKeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigKeyLogic {
	return &ConfigKeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigKeyLogic) ConfigKey(req *types.CodeReq) (resp *types.ConfigKeyResp, err error) {
	q := l.svcCtx.Dal.Query
	config, err := q.SysConfig.WithContext(l.ctx).Where(q.SysConfig.ConfigKey.Eq(req.Code)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = &types.ConfigKeyResp{
		Data: config.ConfigValue,
	}
	return
}
