package data

import (
	"context"
	"github.com/jinzhu/copier"
	"system/internal/svc"
	"system/internal/types"
	"toolkit/errx"

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

func (l *InfoLogic) Info(req *types.IdReq) (resp *types.DictDataBase, err error) {
	resp = new(types.DictDataBase)
	q := l.svcCtx.Dal.Query
	sysDictDatum, err := q.SysDictDatum.WithContext(l.ctx).Where(q.SysDictDatum.DictCode.Eq(req.Id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	if err := copier.Copy(&resp, sysDictDatum); err != nil {
		return nil, err
	}
	return
}
