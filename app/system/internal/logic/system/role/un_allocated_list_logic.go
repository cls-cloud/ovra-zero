package role

import (
	"context"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnAllocatedListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnAllocatedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnAllocatedListLogic {
	return &UnAllocatedListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnAllocatedListLogic) UnAllocatedList(req *types.AllocatedReq) (resp *types.AllocatedResp, err error) {
	list, err := NewAllocatedListLogic(l.ctx, l.svcCtx).GetUserList(req, false)
	if err != nil {
		return nil, err
	}
	resp = new(types.AllocatedResp)
	resp = list
	return
}
