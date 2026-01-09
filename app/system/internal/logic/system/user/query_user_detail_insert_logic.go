package user

import (
	"context"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryUserDetailInsertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryUserDetailInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryUserDetailInsertLogic {
	return &QueryUserDetailInsertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryUserDetailInsertLogic) QueryUserDetailInsert() (resp *types.UserDetailResp, err error) {
	resp = new(types.UserDetailResp)
	userDetail, err := NewQueryUserDetailLogic(l.ctx, l.svcCtx).QueryUserDetail(nil)
	if err != nil {
		return nil, err
	}
	resp = userDetail
	return
}
