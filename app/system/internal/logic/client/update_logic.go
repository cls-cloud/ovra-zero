package client

import (
	"context"
	"ovra/toolkit/errx"
	"ovra/toolkit/utils"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

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

func (l *UpdateLogic) Update(req *types.ModifyClientReq) error {
	toMapOmit := utils.StructToMapOmit(req.ClientBase, nil, nil, true)
	if _, err := l.svcCtx.Dal.Query.SysClient.WithContext(l.ctx).Where(l.svcCtx.Dal.Query.SysClient.ID.Eq(req.ID)).Updates(toMapOmit); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
