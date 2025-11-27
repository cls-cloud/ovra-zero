package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"system/internal/svc"
	"system/internal/types"
	"toolkit/auth"
	"toolkit/errx"
	"toolkit/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResp, err error) {
	// 获取用户权限
	// 获取用户角色
	q := l.svcCtx.Dal.Query
	userId := auth.GetUserId(l.ctx)
	tenantKey := fmt.Sprintf(tenant.TENANT_KEY, userId)
	//删除租户ID缓存
	// 先判断redis中是否存在
	ex, err := l.svcCtx.Rds.ExistsCtx(l.ctx, tenantKey)
	if err != nil {
		return nil, err
	}
	tenantId := auth.GetTenantId(l.ctx)
	if ex {
		tenantId, err = l.svcCtx.Rds.HgetCtx(l.ctx, tenantKey, "ot")
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.Rds.HsetCtx(l.ctx, tenantKey, "nt", tenantId)
		if err != nil {
			return nil, err
		}
	}
	resp = new(types.UserInfoResp)
	sysUser := l.svcCtx.Dal.Query.SysUser
	udo := sysUser.WithContext(l.ctx).Where(sysUser.UserID.Eq(userId))
	if tenantId != "" {
		udo.Where(sysUser.TenantID.Eq(tenantId))
	}
	user, err := udo.First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errx.GORMErrMsg(err, "用户不存在")
		}
		return nil, err
	}
	sysRole := l.svcCtx.Dal.Query.SysRole
	sysUserRole := l.svcCtx.Dal.Query.SysUserRole
	var roleIds []string
	err = sysUserRole.WithContext(l.ctx).Select(sysUserRole.RoleID).Where(sysUserRole.UserID.Eq(userId)).Scan(&roleIds)
	if err != nil {
		return nil, err
	}
	isAdmin := false
	for _, item := range roleIds {
		if item == "1" {
			isAdmin = true
		}
	}
	sysRole.WithContext(l.ctx).Where(sysRole.RoleID.In())
	rdo := sysRole.WithContext(l.ctx).
		Where(sysRole.RoleID.In(roleIds...))
	if tenantId != "" {
		rdo.Where(sysRole.TenantID.Eq(tenantId))
	}
	roles, err := rdo.Find()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errx.GORMErr(err)
		}
		return nil, err
	}
	roleKeys := make([]string, len(roles))
	for i, item := range roles {
		roleKeys[i] = item.RoleKey
	}
	if isAdmin {
		resp.Permissions = append(resp.Permissions, "*:*:*")
		resp.Roles = append(resp.Roles, "superadmin")
	} else {
		resp.Roles = roleKeys
		permissions := make([]string, 0)
		err = q.SysMenu.WithContext(l.ctx).Select(q.SysMenu.Perms).
			LeftJoin(q.SysRoleMenu, q.SysRoleMenu.MenuID.EqCol(q.SysMenu.MenuID)).
			LeftJoin(q.SysUserRole, q.SysUserRole.RoleID.EqCol(q.SysRoleMenu.RoleID)).
			Where(q.SysUserRole.UserID.Eq(user.UserID)).
			Where(q.SysMenu.Perms.IsNotNull(), q.SysMenu.Perms.Neq("")).
			Distinct(q.SysMenu.Perms).
			Scan(&permissions)
		if err != nil {
			return nil, err
		}
		resp.Permissions = permissions
	}

	err = copier.Copy(&resp.User, user)
	err = copier.Copy(&resp.User.Roles, roles)
	if err != nil {
		return nil, err
	}
	return
}
