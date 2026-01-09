package role

import (
	"context"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelectAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSelectAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectAllLogic {
	return &SelectAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelectAllLogic) SelectAll(req *types.SelectAllReq) error {
	userIds := strings.Split(req.UserIds, ",")
	dal := l.svcCtx.Dal
	err := dal.SysRoleDal.AddSysRoleUsers(l.ctx, req.RoleId, userIds)
	if err != nil {
		return err
	}
	return nil
}
