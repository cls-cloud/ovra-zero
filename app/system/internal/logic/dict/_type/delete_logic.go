package _type

import (
	"context"
	"strings"
	"toolkit/errx"

	"system/internal/svc"
	"system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.CodeReq) error {
	ids := strings.Split(req.Code, ",")
	q := l.svcCtx.Dal.Query
	dictTypes := make([]string, 0)
	err := q.SysDictType.WithContext(l.ctx).Select(q.SysDictType.DictType).Where(q.SysDictType.DictID.In(ids...)).Scan(&dictTypes)
	if err != nil {
		return errx.GORMErr(err)
	}
	if _, err := q.SysDictType.WithContext(l.ctx).Where(q.SysDictType.DictID.In(ids...)).Unscoped().Delete(); err != nil {
		return errx.GORMErr(err)
	}
	//删除关联的字典数据
	if _, err := q.SysDictDatum.WithContext(l.ctx).Where(q.SysDictDatum.DictType.In(dictTypes...)).Unscoped().Delete(); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
