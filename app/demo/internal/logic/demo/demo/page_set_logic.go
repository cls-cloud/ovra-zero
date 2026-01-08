// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package demo

import (
	"context"
	"ovra/app/demo/internal/svc"
	"ovra/app/demo/internal/types"
	"strconv"

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

func (l *PageSetLogic) PageSet(req *types.PageSetDemoReq) (resp *types.PageSetDemoResp, err error) {
	resp = new(types.PageSetDemoResp)
	total, list, err := l.svcCtx.Dal.TestDemoDal.PageSet(l.ctx, int(req.PageNum), int(req.PageSize), &req.DemoQuery)
	if err != nil {
		return
	}
	resp.Total = total
	for _, demo := range list {
		resp.Rows = append(resp.Rows, &types.DemoBase{
			ID:       demo.ID,
			TestKey:  demo.TestKey,
			Value:    demo.Value,
			OrderNum: strconv.FormatInt(int64(demo.OrderNum), 10),
			Version:  strconv.FormatInt(int64(demo.Version), 10),
		})
	}
	return
}
