package _post

import (
	"context"
	"fmt"
	"ovra/toolkit/errx"
	"time"

	"github.com/jinzhu/copier"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageSetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPageSetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageSetLogic {
	return &PageSetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageSetLogic) PageSet(req *types.PageSetRoleReq) (resp *types.PageSetRoleResp, err error) {
	offset := (req.PageNum - 1) * req.PageSize
	q := l.svcCtx.Dal.Query
	do := q.SysPost.WithContext(l.ctx)
	if req.PostName != "" {
		do = do.Where(q.SysPost.PostName.Like(fmt.Sprintf("%%%s%%", req.PostName)))
	}
	if req.PostCode != "" {
		do = do.Where(q.SysPost.PostCode.Like(fmt.Sprintf("%%%s%%", req.PostCode)))
	}
	if req.Status != "" {
		do = do.Where(q.SysPost.Status.Eq(req.Status))
	}
	if req.BelongDeptId != "" {
		do = do.Where(q.SysPost.DeptID.Eq(req.BelongDeptId))
	}
	result, count, err := do.Order(q.SysPost.PostSort.Asc(), q.SysPost.CreateTime.Desc()).FindByPage(int(offset), int(req.PageSize))
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = new(types.PageSetRoleResp)
	resp.Total = count
	list := make([]*types.PostBase, len(result))
	for i, item := range result {
		list[i] = new(types.PostBase)
		if err = copier.Copy(&list[i], item); err != nil {
			return nil, err
		}
		list[i].CreateTime = item.CreateTime.Format(time.DateTime)
	}
	resp.Rows = list
	return
}
