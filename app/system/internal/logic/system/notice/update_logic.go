package notice

import (
	"context"
	"ovra/app/system/internal/dal/model"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

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

func (l *UpdateLogic) Update(req *types.ModifyNoticeReq) error {
	if err := l.svcCtx.Dal.SysNoticeDal.Update(l.ctx, &model.SysNotice{
		NoticeID:      req.NoticeID,
		NoticeTitle:   req.NoticeTitle,
		NoticeType:    req.NoticeType,
		NoticeContent: []byte(req.NoticeContent),
		Status:        req.Status,
		Remark:        req.Remark,
	}); err != nil {
		return err
	}
	return nil
}
