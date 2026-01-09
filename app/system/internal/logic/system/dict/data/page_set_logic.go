package data

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

func (l *PageSetLogic) PageSet(req *types.PageSetDictDataReq) (resp *types.PageSetDictDataResp, err error) {
	offset := (req.PageNum - 1) * req.PageSize
	q := l.svcCtx.Dal.Query
	do := q.SysDictDatum.WithContext(l.ctx)
	if req.DictType != "" {
		do = do.Where(q.SysDictDatum.DictType.Like(fmt.Sprintf("%%%s%%", req.DictType)))
	}

	result, count, err := do.Order(q.SysDictDatum.DictSort.Asc(), q.SysDictDatum.CreateTime.Desc()).FindByPage(int(offset), int(req.PageSize))
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = new(types.PageSetDictDataResp)
	resp.Total = count
	list := make([]*types.DictDataBase, len(result))
	for i, item := range result {
		list[i] = new(types.DictDataBase)
		if err = copier.Copy(&list[i], item); err != nil {
			return nil, err
		}
		list[i].CreateTime = item.CreateTime.Format(time.DateTime)
	}
	resp.Rows = list
	return
}
