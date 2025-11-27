package _type

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"time"
	"toolkit/errx"

	"system/internal/svc"
	"system/internal/types"

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

func (l *PageSetLogic) PageSet(req *types.PageSetDictTypeReq) (resp *types.PageSetDictTypeResp, err error) {
	offset := (req.PageNum - 1) * req.PageSize
	q := l.svcCtx.Dal.Query
	do := q.SysDictType.WithContext(l.ctx)
	if req.DictName != "" {
		do = do.Where(q.SysDictType.DictName.Like(fmt.Sprintf("%%%s%%", req.DictName)))
	}
	if req.DictType != "" {
		do = do.Where(q.SysDictType.DictType.Like(fmt.Sprintf("%%%s%%", req.DictType)))
	}

	result, count, err := do.Order(q.SysDictType.CreateTime.Desc()).FindByPage(int(offset), int(req.PageSize))
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = new(types.PageSetDictTypeResp)
	resp.Total = count
	list := make([]*types.DictTypeBase, len(result))
	for i, item := range result {
		list[i] = new(types.DictTypeBase)
		if err = copier.Copy(&list[i], item); err != nil {
			return nil, err
		}
		list[i].CreateTime = item.CreateTime.Format(time.DateTime)
	}
	resp.Rows = list
	return
}
