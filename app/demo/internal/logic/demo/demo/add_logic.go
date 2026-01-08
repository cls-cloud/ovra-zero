// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package demo

import (
	"context"
	"ovra/app/demo/internal/dal/model"
	"ovra/app/demo/internal/svc"
	"ovra/app/demo/internal/types"

	"github.com/ovra-cloud/ovra-toolkit/auth"
	"github.com/ovra-cloud/ovra-toolkit/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req *types.ModifyDemoReq) error {
	userId := auth.GetUserId(l.ctx)
	param := &model.TestDemo{
		UserID:   utils.StrAtoi(userId),
		OrderNum: int32(utils.StrAtoi(req.OrderNum)),
		TestKey:  req.TestKey,
		Value:    req.Value,
		Version:  int32(utils.StrAtoi(req.Version)),
	}
	return l.svcCtx.Dal.TestDemoDal.Insert(l.ctx, param)
}
