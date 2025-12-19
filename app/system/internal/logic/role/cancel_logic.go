package role

import (
	"context"
	"ovra/toolkit/errx"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLogic {
	return &CancelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelLogic) Cancel(req *types.CancelReq) error {
	q := l.svcCtx.Dal.Query
	roleId := req.RoleId
	if _, err := q.SysUserRole.WithContext(l.ctx).Where(q.SysUserRole.UserID.Eq(req.UserId)).Where(q.SysUserRole.RoleID.Eq(roleId)).Unscoped().Delete(); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
