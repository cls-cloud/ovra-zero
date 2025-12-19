package role

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

func (l *ChangeStatusLogic) ChangeStatus(req *types.UpdateRoleStatusReq) error {
	q := l.svcCtx.Dal.Query
	_, err := q.SysRole.WithContext(l.ctx).Where(q.SysRole.RoleID.Eq(req.RoleID)).Update(q.SysRole.Status, req.Status)
	if err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
