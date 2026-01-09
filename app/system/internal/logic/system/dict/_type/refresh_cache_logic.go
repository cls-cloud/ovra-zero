package _type

import (
	"context"
	"encoding/json"
	"ovra/toolkit/constants"
	"ovra/toolkit/errx"

	"ovra/app/system/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshCacheLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshCacheLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshCacheLogic {
	return &RefreshCacheLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshCacheLogic) RefreshCache() error {
	q := l.svcCtx.Dal.Query
	key := constants.DictCache
	dictTypes, err := q.SysDictType.WithContext(l.ctx).Find()
	if err != nil {
		return errx.GORMErr(err)
	}
	_, err = l.svcCtx.Rds.DelCtx(l.ctx, key)
	if err != nil {
		return err
	}
	for _, dictType := range dictTypes {
		dictDatum, err := q.SysDictDatum.WithContext(l.ctx).Where(q.SysDictDatum.DictType.Eq(dictType.DictType)).Find()
		if err != nil {
			return errx.GORMErr(err)
		}
		bytes, err := json.Marshal(dictDatum)
		err = l.svcCtx.Rds.HsetCtx(l.ctx, key, dictType.DictType, string(bytes))
		if err != nil {
			return err
		}
	}
	return nil
}
