// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package demo

import (
	"context"
	"ovra/app/demo/internal/svc"
	"ovra/app/demo/internal/types"
	"strconv"

	"github.com/ovra-cloud/ovra-toolkit/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoLogic) Info(req *types.IdReq) (resp *types.DemoBase, err error) {
	data, err := l.svcCtx.Dal.TestDemoDal.SelectById(l.ctx, utils.StrAtoi(req.Id))
	if err != nil {
		return nil, err
	}
	resp = &types.DemoBase{
		ID:       data.ID,
		TestKey:  data.TestKey,
		Value:    data.Value,
		OrderNum: strconv.FormatInt(int64(data.OrderNum), 10),
		Version:  strconv.FormatInt(int64(data.Version), 10),
	}
	return
}
