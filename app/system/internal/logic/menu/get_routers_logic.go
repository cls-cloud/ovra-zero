package menu

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"system/internal/dal/model"
	"system/internal/svc"
	"system/internal/types"
	"toolkit/auth"
)

type GetRoutersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoutersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoutersLogic {
	return &GetRoutersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoutersLogic) GetRouters() (resp []*types.RouterMenuResp, err error) {
	resp = make([]*types.RouterMenuResp, 0)
	userId := auth.GetUserId(l.ctx)
	// 判断用户是否属于超级管理员
	isAdmin := false
	sysUserRole := l.svcCtx.Dal.Query.SysUserRole
	role, err := sysUserRole.WithContext(l.ctx).Where(sysUserRole.UserID.Eq(userId)).Where(sysUserRole.RoleID.Eq("1")).First()
	if err == nil && role != nil {
		isAdmin = true
	}
	sysMenus, err := l.GetMenuByUserId(l.ctx, userId, isAdmin)
	if err != nil {
		return
	}
	var menus []*types.RouterMenuResp
	for idx, bizMenu := range sysMenus {
		isRoot := bizMenu.ParentID == "0"
		isHttpLink := strings.HasPrefix(bizMenu.Path, "http://") || strings.HasPrefix(bizMenu.Path, "https://")

		router := &types.RouterMenuResp{
			MenuId:   bizMenu.MenuID,
			ParentId: bizMenu.ParentID,
			Hidden:   bizMenu.Visible == "1",
			Name:     fmt.Sprintf("%s.%d", strings.Title(bizMenu.Path), idx),
			Meta: &types.RouterMenuMeta{
				Title:   bizMenu.MenuName,
				NoCache: bizMenu.IsCache == 1,
				Icon:    bizMenu.Icon,
			},
		}

		// 设置 Path、Redirect、Component、AlwaysShow
		switch {
		case bizMenu.IsFrame == External && isHttpLink:
			// 外链跳转
			router.Path = bizMenu.Path
			router.Redirect = bizMenu.Path
			router.Component = Layout

		case isHttpLink:
			// 内嵌 iframe
			router.Path = bizMenu.Path
			router.Redirect = bizMenu.Path
			router.Component = InnerLink

		case isRoot:
			// 一级菜单
			router.Path = "/" + bizMenu.Path
			router.Redirect = "noRedirect"
			router.Component = Layout
			router.AlwaysShow = true

		default:
			// 默认子菜单
			router.Path = bizMenu.Path
			router.Redirect = ""
			if bizMenu.Component == "" {
				bizMenu.Component = ParentView
			}
			router.Component = bizMenu.Component
		}

		menus = append(menus, router)
	}

	resp = l.Tree(menus, "0")
	return
}

func (l *GetRoutersLogic) Tree(node []*types.RouterMenuResp, pid string) []*types.RouterMenuResp {
	res := make([]*types.RouterMenuResp, 0)
	for _, v := range node {
		if v.ParentId == pid {
			v.Children = l.Tree(node, v.MenuId)
			res = append(res, v)
		}
	}
	return res
}
func (l *GetRoutersLogic) GetMenuByUserId(ctx context.Context, userId string, isAdmin bool) ([]*model.SysMenu, error) {
	q := l.svcCtx.Dal.Query

	var sysMenus []*model.SysMenu
	var err error
	if isAdmin {
		if sysMenus, err = q.SysMenu.WithContext(ctx).
			Where(q.SysMenu.MenuType.In(TypeDir, TypeMenu)).
			Where(q.SysMenu.Status.Eq("0")).
			Order(q.SysMenu.ParentID, q.SysMenu.OrderNum).
			Find(); err != nil {
			return nil, err
		}
	} else {
		if sysMenus, err = q.SysMenu.WithContext(ctx).
			Distinct().
			LeftJoin(q.SysRoleMenu, q.SysRoleMenu.MenuID.EqCol(q.SysMenu.MenuID)).
			LeftJoin(q.SysUserRole, q.SysRoleMenu.RoleID.EqCol(q.SysUserRole.RoleID)).
			LeftJoin(q.SysRole, q.SysUserRole.RoleID.EqCol(q.SysRole.RoleID)).
			LeftJoin(q.SysUser, q.SysUserRole.UserID.EqCol(q.SysUser.UserID)).
			Where(
				q.SysUser.UserID.Eq(userId),
				q.SysMenu.MenuType.In(TypeDir, TypeMenu),
				q.SysMenu.Status.Eq("0"),
				q.SysRole.Status.Eq("0"),
			).
			Order(q.SysMenu.ParentID, q.SysMenu.OrderNum).
			Find(); err != nil {
			return nil, err
		}
	}
	return sysMenus, nil
}
