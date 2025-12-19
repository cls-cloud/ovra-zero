package data

import (
	"context"
	"ovra/toolkit/errx"
	"strings"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDictDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictDataLogic {
	return &DeleteDictDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDictDataLogic) DeleteDictData(req *types.CodeReq) error {
	ids := strings.Split(req.Code, ",")
	q := l.svcCtx.Dal.Query
	if _, err := q.SysDictDatum.WithContext(l.ctx).Where(q.SysDictDatum.DictCode.In(ids...)).Unscoped().Delete(); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
