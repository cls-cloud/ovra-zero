package menu

import (
	"context"
	"fmt"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"ovra/toolkit/auth"
	"ovra/toolkit/errx"
	"time"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.QueryMenuListReq) (resp []*types.MenuBase, err error) {
	q := l.svcCtx.Dal.Query
	userId := auth.GetUserId(l.ctx)
	do := q.SysMenu.WithContext(l.ctx)
	// 除了超级管理员 其余所有用户都只能查询自己角色绑定的菜单
	if userId != "1" {
		do = do.LeftJoin(q.SysRoleMenu, q.SysRoleMenu.MenuID.EqCol(q.SysMenu.MenuID)).
			LeftJoin(q.SysUserRole, q.SysUserRole.RoleID.EqCol(q.SysRoleMenu.RoleID)).
			Where(q.SysUserRole.UserID.Eq(userId))
	}
	if req.MenuName != "" {
		do = do.Where(q.SysMenu.MenuName.Like(fmt.Sprintf("%%%s%%", req.MenuName)))
	}
	if req.Status != "" {
		do = do.Where(q.SysMenu.Status.Eq(req.Status))
	}
	if req.Visible != "" {
		do = do.Where(q.SysMenu.Visible.Eq(req.Visible))
	}

	result, err := do.Order(q.SysMenu.OrderNum.Asc(), q.SysMenu.CreateTime.Desc()).Find()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = make([]*types.MenuBase, 0)
	list := make([]*types.MenuBase, len(result))
	for i, item := range result {
		list[i] = new(types.MenuBase)
		if err = copier.Copy(&list[i], item); err != nil {
			return nil, err
		}
		list[i].CreateTime = item.CreateTime.Format(time.DateTime)
	}
	resp = list
	return
}
