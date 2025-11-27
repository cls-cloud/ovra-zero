package menu

import (
	"context"
	"strings"
	"toolkit/errx"

	"system/internal/svc"
	"system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantPackageTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTenantPackageTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantPackageTreeLogic {
	return &TenantPackageTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *TenantPackageTreeLogic) TenantPackageTree(req *types.IdReq) (resp *types.SelectMenuTreeResp, err error) {
	resp = new(types.SelectMenuTreeResp)
	q := l.svcCtx.Dal.Query
	resultIds := make([]string, 0)

	// 只有非 0 ID 才查套餐
	if req.Id != "0" {
		// 查询租户套餐
		tenantPackage, err := q.SysTenantPackage.WithContext(l.ctx).
			Where(q.SysTenantPackage.PackageID.Eq(req.Id)).First()
		if err != nil {
			return nil, errx.GORMErr(err)
		}

		// 解析 menuIds
		var menuIds []string
		if tenantPackage.MenuIds != "" {
			menuIds = strings.Split(tenantPackage.MenuIds, ",")
		}

		if len(menuIds) > 0 {
			// 若为关联模式，查出 parentId 列表
			var parentIds []string
			if tenantPackage.MenuCheckStrictly {
				err := q.SysMenu.WithContext(l.ctx).
					Select(q.SysMenu.ParentID).
					Where(q.SysMenu.MenuID.In(menuIds...)).
					Distinct(q.SysMenu.ParentID).
					Scan(&parentIds)
				if err != nil {
					return nil, errx.GORMErr(err)
				}
			}

			// 获取最终的 menuId（排除父菜单）
			query := q.SysMenu.WithContext(l.ctx).
				Select(q.SysMenu.MenuID).
				Where(q.SysMenu.MenuID.In(menuIds...))
			if len(parentIds) > 0 {
				query = query.Where(q.SysMenu.MenuID.NotIn(parentIds...))
			}

			if err := query.Scan(&resultIds); err != nil {
				return nil, errx.GORMErr(err)
			}
		}
	}

	resp.CheckedKeys = resultIds

	// 获取菜单树结构
	treeSelect, err := NewTreeSelectLogic(l.ctx, l.svcCtx).TreeSelect()
	if err != nil {
		return nil, err
	}
	resp.Menus = treeSelect
	return
}
