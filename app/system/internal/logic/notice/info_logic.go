package notice

import (
	"context"
	"ovra/toolkit/errx"

	"github.com/jinzhu/copier"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoLogic) Info(req *types.IdReq) (resp *types.NoticeBase, err error) {
	resp = new(types.NoticeBase)
	q := l.svcCtx.Dal.Query
	sysNotice, err := q.SysNotice.WithContext(l.ctx).Where(q.SysNotice.NoticeID.Eq(req.Id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	err = copier.Copy(&resp, sysNotice)
	return
}
