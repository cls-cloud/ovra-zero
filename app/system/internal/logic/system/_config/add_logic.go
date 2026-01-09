package _config

import (
	"context"
	"ovra/app/system/internal/dal/model"
	"ovra/toolkit/errx"
	"ovra/toolkit/utils"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req *types.ModifyConfigReq) error {
	config := model.SysConfig{
		ConfigID:    utils.GetID(),
		ConfigKey:   req.ConfigKey,
		ConfigName:  req.ConfigName,
		ConfigType:  req.ConfigType,
		ConfigValue: req.ConfigValue,
		Remark:      req.Remark,
	}
	if err := l.svcCtx.Dal.Query.SysConfig.WithContext(l.ctx).Create(&config); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
