package _config

import (
	"context"
	"toolkit/errx"
	"toolkit/utils"

	"system/internal/svc"
	"system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.ModifyConfigReq) error {
	toMapOmit := utils.StructToMapOmit(req.ConfigBase, nil, []string{"CreateTime"}, true)
	if _, err := l.svcCtx.Dal.Query.SysConfig.WithContext(l.ctx).Where(l.svcCtx.Dal.Query.SysConfig.ConfigID.Eq(req.ConfigID)).Updates(toMapOmit); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
