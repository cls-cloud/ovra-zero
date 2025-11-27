package role

import (
	"context"
	"net/http"
	"strings"
	"toolkit/errx"

	"system/internal/svc"
	"system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveLogic {
	return &RemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveLogic) Remove(req *types.CodeReq) error {
	ids := strings.Split(req.Code, ",")
	q := l.svcCtx.Dal.Query
	//先查询角色是否有用户绑定
	if count, err := q.SysUserRole.WithContext(l.ctx).Where(q.SysUserRole.RoleID.In(ids...)).Count(); err != nil {
		return errx.GORMErr(err)
	} else if count > 0 {
		return errx.New(http.StatusForbidden, "角色下有用户绑定，请先解绑用户")
	}
	if _, err := q.SysRole.WithContext(l.ctx).Where(q.SysRole.RoleID.In(ids...)).Unscoped().Delete(); err != nil {
		return errx.GORMErr(err)
	}
	if _, err := q.SysUserRole.WithContext(l.ctx).Where(q.SysUserRole.RoleID.In(ids...)).Unscoped().Delete(); err != nil {
		return errx.GORMErr(err)
	}
	if _, err := q.SysRoleMenu.WithContext(l.ctx).Where(q.SysRoleMenu.RoleID.In(ids...)).Unscoped().Delete(); err != nil {
		return errx.GORMErr(err)
	}
	if _, err := q.SysRoleDept.WithContext(l.ctx).Where(q.SysRoleDept.RoleID.In(ids...)).Unscoped().Delete(); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
