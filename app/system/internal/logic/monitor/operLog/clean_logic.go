package operLog

import (
	"context"
	"ovra/app/system/internal/dal/model"

	"ovra/app/system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CleanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCleanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CleanLogic {
	return &CleanLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CleanLogic) Clean() error {
	q := l.svcCtx.Dal.Query
	_, err := q.SysOperLog.WithContext(l.ctx).Unscoped().Delete(&model.SysOperLog{})
	if err != nil {
		return err
	}
	return nil
}
