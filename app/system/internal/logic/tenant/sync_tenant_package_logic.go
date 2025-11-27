package tenant

import (
	"context"
	"strings"
	"system/internal/dal/model"
	"toolkit/errx"

	"system/internal/svc"
	"system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncTenantPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncTenantPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncTenantPackageLogic {
	return &SyncTenantPackageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *SyncTenantPackageLogic) SyncTenantPackage(req *types.SyncTenantPackageReq) error {
	q := l.svcCtx.Dal.Query
	// 获取套餐信息
	tenantPackage, err := q.SysTenantPackage.
		WithContext(l.ctx).
		Where(q.SysTenantPackage.PackageID.Eq(req.PackageId)).
		First()
	if err != nil {
		return errx.GORMErr(err)
	}

	// 查询当前租户下所有角色
	roles, err := q.SysRole.
		WithContext(l.ctx).
		Where(q.SysRole.TenantID.Eq(req.TenantId)).
		Find()
	if err != nil {
		return errx.GORMErr(err)
	}
	// 提取 menuIds
	menuIds := strings.Split(tenantPackage.MenuIds, ",")
	var roleIds []string
	for _, role := range roles {
		if role.RoleKey == "admin" {
			// 管理员角色重新绑定菜单
			var roleMenus []*model.SysRoleMenu
			for _, menuId := range menuIds {
				roleMenus = append(roleMenus, &model.SysRoleMenu{
					RoleID: role.RoleID,
					MenuID: menuId,
				})
			}
			// 删除旧绑定
			_, err := q.SysRoleMenu.
				WithContext(l.ctx).
				Where(q.SysRoleMenu.RoleID.Eq(role.RoleID)).
				Unscoped().
				Delete()
			if err != nil {
				return err
			}
			// 批量插入新绑定
			if len(roleMenus) > 0 {
				err := q.SysRoleMenu.
					WithContext(l.ctx).
					CreateInBatches(roleMenus, len(roleMenus))
				if err != nil {
					return err
				}
			}
		} else {
			roleIds = append(roleIds, role.RoleID)
		}
	}
	// 非管理员角色，清除不在套餐中的菜单
	if len(roleIds) > 0 {
		_, err := q.SysRoleMenu.
			WithContext(l.ctx).
			Where(
				q.SysRoleMenu.RoleID.In(roleIds...),
				q.SysRoleMenu.MenuID.NotIn(menuIds...),
			).
			Unscoped().
			Delete()
		if err != nil {
			return err
		}
	}
	return nil
}
