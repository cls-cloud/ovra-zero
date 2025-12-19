package dept

import (
	"context"
	"fmt"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
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

func (l *ListLogic) List(req *types.DeptQuery) (resp []*types.DeptBase, err error) {
	q := l.svcCtx.Dal.Query
	do := q.SysDept.WithContext(l.ctx)
	if req.DeptName != "" {
		do = do.Where(q.SysDept.DeptName.Like(fmt.Sprintf("%%%s%%", req.DeptName)))
	}
	if req.Status != "" {
		do = do.Where(q.SysDept.Status.Eq(req.Status))
	}
	result, err := do.Order(q.SysDept.OrderNum.Asc(), q.SysDept.CreateTime.Desc()).Find()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = make([]*types.DeptBase, 0)
	list := make([]*types.DeptBase, len(result))
	for i, item := range result {
		list[i] = new(types.DeptBase)
		if err = copier.Copy(&list[i], item); err != nil {
			return nil, err
		}
		list[i].DeptID = item.DeptID
		list[i].ParentID = item.ParentID
		list[i].CreateTime = item.CreateTime.Format(time.DateTime)
	}
	resp = list
	return
}
