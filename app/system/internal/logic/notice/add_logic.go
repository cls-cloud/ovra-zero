package notice

import (
	"context"
	"ovra/app/system/internal/dal/model"
	"ovra/toolkit/errx"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

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

func (l *AddLogic) Add(req *types.ModifyNoticeReq) error {
	notice := &model.SysNotice{
		NoticeTitle:   req.NoticeTitle,
		NoticeType:    req.NoticeType,
		NoticeContent: []byte(req.NoticeContent),
		Status:        req.Status,
	}
	if err := l.svcCtx.Dal.Query.SysNotice.WithContext(l.ctx).Create(notice); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
