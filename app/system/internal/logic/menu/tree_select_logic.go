package menu

import (
	"context"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"ovra/toolkit/auth"
	"ovra/toolkit/errx"

	"github.com/zeromicro/go-zero/core/logx"
)

type TreeSelectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTreeSelectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TreeSelectLogic {
	return &TreeSelectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TreeSelectLogic) TreeSelect() (resp []*types.SelectMenuTree, err error) {
	resp = make([]*types.SelectMenuTree, 0)
	q := l.svcCtx.Dal.Query

	userId := auth.GetUserId(l.ctx)
	do := q.SysMenu.WithContext(l.ctx)
	// 除了超级管理员 其余所有用户都只能查询自己角色绑定的菜单
	if userId != "1" {
		do = do.LeftJoin(q.SysRoleMenu, q.SysRoleMenu.MenuID.EqCol(q.SysMenu.MenuID)).
			LeftJoin(q.SysUserRole, q.SysUserRole.RoleID.EqCol(q.SysRoleMenu.RoleID)).
			Where(q.SysUserRole.UserID.Eq(userId))
	}
	sysMenus, err := do.Order(q.SysMenu.OrderNum.Asc(), q.SysMenu.CreateTime.Desc()).Find()

	if err != nil {
		return nil, errx.GORMErr(err)
	}
	var menuTree []*types.SelectMenuTree
	for _, menu := range sysMenus {
		menuTree = append(menuTree, &types.SelectMenuTree{
			Id:       menu.MenuID,
			ParentId: menu.ParentID,
			Icon:     menu.Icon,
			MenuType: menu.MenuType,
			Label:    menu.MenuName,
		})
	}
	tree := l.Tree(menuTree, "0")
	resp = tree

	return
}

func (l *TreeSelectLogic) Tree(node []*types.SelectMenuTree, pid string) []*types.SelectMenuTree {
	res := make([]*types.SelectMenuTree, 0)
	for _, v := range node {
		if v.ParentId == pid {
			v.Children = l.Tree(node, v.Id)
			res = append(res, v)
		}
	}
	return res
}
