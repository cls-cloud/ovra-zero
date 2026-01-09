package tenant

import (
	"context"
	"ovra/toolkit/auth"
	"ovra/toolkit/tenant"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DynamicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDynamicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DynamicLogic {
	return &DynamicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DynamicLogic) Dynamic(req *types.IdReq) error {
	userId := auth.GetUserId(l.ctx)
	if err := tenant.SetTenantId(l.ctx, l.svcCtx.Rds, userId, req.Id); err != nil {
		return err
	}
	return nil
}
