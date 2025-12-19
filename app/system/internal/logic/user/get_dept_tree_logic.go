package user

import (
	"context"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"ovra/toolkit/errx"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeptTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeptTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeptTreeLogic {
	return &GetDeptTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeptTreeLogic) GetDeptTree() (resp []types.DeptTree, err error) {
	q := l.svcCtx.Dal.Query
	sysDepts, err := q.SysDept.WithContext(l.ctx).Order(q.SysDept.OrderNum).Find()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	for _, dept := range sysDepts {
		resp = append(resp, types.DeptTree{
			Id:       dept.DeptID,
			ParentId: dept.ParentID,
			Label:    dept.DeptName,
		})
	}
	resp = BuildDeptTree(resp, "0")
	return
}

func BuildDeptTree(list []types.DeptTree, pid string) []types.DeptTree {
	var tree []types.DeptTree
	for _, item := range list {
		if item.ParentId == pid {
			children := BuildDeptTree(list, item.Id)
			if len(children) > 0 {
				item.Children = children
			}
			tree = append(tree, item)
		}
	}
	return tree
}
