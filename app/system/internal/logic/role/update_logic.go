package role

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

func (l *UpdateLogic) Update(req *types.AddOrUpdateRoleReq) error {
	if req.RoleID == "" {
		return errx.BizErr("角色ID不能为空")
	}
	roleId := req.RoleID
	q := l.svcCtx.Query
	dal := l.svcCtx.Dal
	if exit := dal.SysRoleDal.SelectByRoleKeyExit(l.ctx, roleId, req.RoleKey); exit {
		return errx.BizErr("角色编码已存在")
	}

	omit := utils.StructToMapOmit(req.RoleBase, nil, []string{"SuperAdmin"}, true)
	if _, err := q.SysRole.WithContext(l.ctx).Where(q.SysRole.RoleID.Eq(roleId)).Updates(omit); err != nil {
		return errx.GORMErr(err)
	}
	if len(req.MenuIds) != 0 {
		if err := dal.SysRoleDal.AddSysRoleMenus(l.ctx, roleId, req.MenuIds); err != nil {
			return err
		}
	}
	return nil
}
