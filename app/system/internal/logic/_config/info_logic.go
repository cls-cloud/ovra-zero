package _config

import (
	"context"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"ovra/toolkit/errx"

	"github.com/jinzhu/copier"

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

func (l *InfoLogic) Info(req *types.IdReq) (resp *types.ConfigBase, err error) {
	resp = new(types.ConfigBase)
	q := l.svcCtx.Dal.Query
	sysPost, err := q.SysConfig.WithContext(l.ctx).Where(q.SysConfig.ConfigID.Eq(req.Id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	err = copier.Copy(&resp, sysPost)
	return
}
