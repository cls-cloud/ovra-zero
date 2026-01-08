// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package demo

import (
	"context"
	"ovra/app/demo/internal/svc"
	"ovra/app/demo/internal/types"

	"github.com/ovra-cloud/ovra-toolkit/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.IdsReq) error {
	ids, _ := utils.SplitToInt64(req.Ids, ",")
	return l.svcCtx.Dal.TestDemoDal.DeleteBatch(l.ctx, ids)
}
