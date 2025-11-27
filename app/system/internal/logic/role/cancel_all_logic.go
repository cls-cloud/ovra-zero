package role

import (
	"context"
	"strings"
	"toolkit/errx"

	"system/internal/svc"
	"system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelAllLogic {
	return &CancelAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelAllLogic) CancelAll(req *types.SelectAllReq) error {
	userIds := strings.Split(req.UserIds, ",")
	q := l.svcCtx.Dal.Query
	if _, err := q.SysUserRole.WithContext(l.ctx).Where(q.SysUserRole.UserID.In(userIds...)).Unscoped().Delete(); err != nil {
		return errx.GORMErr(err)
	}

	return nil
}
