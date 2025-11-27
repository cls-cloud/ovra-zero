package _config

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

func (l *PageSetLogic) PageSet(req *types.PageSetConfigReq) (resp *types.PageSetConfigResp, err error) {
	offset := (req.PageNum - 1) * req.PageSize
	q := l.svcCtx.Dal.Query
	do := q.SysConfig.WithContext(l.ctx)
	if req.ConfigName != "" {
		do = do.Where(q.SysConfig.ConfigName.Like(fmt.Sprintf("%%%s%%", req.ConfigName)))
	}
	if req.ConfigKey != "" {
		do = do.Where(q.SysConfig.ConfigKey.Like(fmt.Sprintf("%%%s%%", req.ConfigKey)))
	}
	result, count, err := do.Order(q.SysConfig.CreateTime.Desc()).FindByPage(int(offset), int(req.PageSize))
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = new(types.PageSetConfigResp)
	resp.Total = count
	list := make([]*types.ConfigBase, len(result))
	for i, item := range result {
		list[i] = new(types.ConfigBase)
		if err = copier.Copy(&list[i], item); err != nil {
			return nil, err
		}
		list[i].CreateTime = item.CreateTime.Format(time.DateTime)
	}
	resp.Rows = list
	return
}
