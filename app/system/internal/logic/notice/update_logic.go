package notice

import (
	"context"
	"toolkit/errx"
	"toolkit/utils"

	"system/internal/svc"
	"system/internal/types"

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
	toMapOmit := utils.StructToMapOmit(req.NoticeBase, nil, []string{"CreateTime"}, true)
	if _, err := l.svcCtx.Dal.Query.SysNotice.WithContext(l.ctx).Where(l.svcCtx.Dal.Query.SysNotice.NoticeID.Eq(req.NoticeID)).Updates(toMapOmit); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
