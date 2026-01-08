// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package demo

import (
	"context"
	"ovra/app/demo/internal/dal/model"
	"ovra/app/demo/internal/svc"
	"ovra/app/demo/internal/types"

	"github.com/ovra-cloud/ovra-toolkit/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.ModifyDemoReq) error {
	if err := l.svcCtx.Dal.TestDemoDal.Update(l.ctx, &model.TestDemo{
		ID:       req.ID,
		OrderNum: int32(utils.StrAtoi(req.OrderNum)),
		TestKey:  req.TestKey,
		Value:    req.Value,
		Version:  int32(utils.StrAtoi(req.Version)),
	}); err != nil {
		return err
	}
	return nil
}
