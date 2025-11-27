package role

import (
	"context"
	"system/internal/dal/model"
	"system/internal/svc"
	"system/internal/types"
	"toolkit/errx"
	"toolkit/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req *types.AddOrUpdateRoleReq) error {
	dal := l.svcCtx.Dal
	if exit := dal.SysRoleDal.SelectByRoleKeyExit(l.ctx, "", req.RoleKey); exit {
		return errx.BizErr("角色编码已存在")
	}
	roleId := utils.GetID()
	if req.RoleID != "" {
		roleId = req.RoleID
	}
	role := &model.SysRole{
		RoleID:    roleId,
		RoleKey:   req.RoleKey,
		RoleName:  req.RoleName,
		RoleSort:  req.RoleSort,
		Status:    req.Status,
		Remark:    req.Remark,
		DataScope: req.DataScope,
	}
	if req.TenantID != "" {
		role.TenantID = req.TenantID
	}
	if err := dal.SysRoleDal.Insert(l.ctx, role); err != nil {
		return err
	}
	if len(req.MenuIds) != 0 {
		if err := dal.SysRoleDal.AddSysRoleMenus(l.ctx, roleId, req.MenuIds); err != nil {
			return err
		}
	}
	return nil
}
