package tenant

import (
	"context"
	"ovra/toolkit/errx"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeStatusLogic {
	return &ChangeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeStatusLogic) ChangeStatus(req *types.ChangeStatusTenantReq) error {
	q := l.svcCtx.Dal.Query
	_, err := q.SysTenant.WithContext(l.ctx).Where(q.SysTenant.TenantID.Eq(req.TenantId), q.SysTenant.ID.Eq(req.Id)).
		Update(q.SysTenant.Status, req.Status)
	if err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
