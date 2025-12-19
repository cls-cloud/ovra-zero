package userrpclogic

import (
	"context"

	"ovra/app/system/internal/svc"
	"ovra/app/system/pb/system"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *system.Request) (*system.Response, error) {
	// todo: add your logic here and delete this line

	return &system.Response{}, nil
}
