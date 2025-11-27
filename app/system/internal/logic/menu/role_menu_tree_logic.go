package menu

import (
	"context"
	"system/internal/svc"
	"system/internal/types"
	"toolkit/errx"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleMenuTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleMenuTreeLogic {
	return &RoleMenuTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *RoleMenuTreeLogic) RoleMenuTree(req *types.IdReq) (resp *types.SelectMenuTreeResp, err error) {
	resp = new(types.SelectMenuTreeResp)
	q := l.svcCtx.Dal.Query

	// 查询角色信息
	role, err := q.SysRole.WithContext(l.ctx).Where(q.SysRole.RoleID.Eq(req.Id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}

	// 查询角色菜单
	roleMenus, err := q.SysRoleMenu.WithContext(l.ctx).
		LeftJoin(q.SysMenu, q.SysMenu.MenuID.EqCol(q.SysRoleMenu.MenuID)).
		Where(q.SysRoleMenu.RoleID.Eq(req.Id)).
		Order(q.SysMenu.ParentID, q.SysMenu.OrderNum).
		Find()
	if err != nil {
		return nil, errx.GORMErr(err)
	}

	menuIds := make([]string, 0, len(roleMenus))
	for _, rm := range roleMenus {
		menuIds = append(menuIds, rm.MenuID)
	}

	// 处理 menuCheckStrictly 模式
	if len(menuIds) != 0 && role.MenuCheckStrictly {
		// 获取父 ID
		var parentIds []string
		err := q.SysMenu.WithContext(l.ctx).
			Select(q.SysMenu.ParentID).
			Join(q.SysRoleMenu, q.SysMenu.MenuID.EqCol(q.SysRoleMenu.MenuID)).
			Where(q.SysRoleMenu.RoleID.Eq(req.Id)).
			Scan(&parentIds)
		if err != nil {
			return nil, errx.GORMErr(err)
		}

		// 去重 & 排除 parentID
		parentIdSet := make(map[string]struct{})
		for _, pid := range parentIds {
			parentIdSet[pid] = struct{}{}
		}

		filteredMenuIds := make([]string, 0)
		for _, id := range menuIds {
			if _, isParent := parentIdSet[id]; !isParent {
				filteredMenuIds = append(filteredMenuIds, id)
			}
		}
		resp.CheckedKeys = filteredMenuIds
	} else {
		resp.CheckedKeys = menuIds
	}

	// 构建菜单树
	treeSelect, err := NewTreeSelectLogic(l.ctx, l.svcCtx).TreeSelect()
	if err != nil {
		return nil, err
	}
	resp.Menus = treeSelect
	return
}
