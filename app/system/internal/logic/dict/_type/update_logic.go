package _type

import (
	"context"
	"ovra/toolkit/errx"
	"ovra/toolkit/utils"

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

func (l *UpdateLogic) Update(req *types.ModifyDictTypeReq) error {
	toMapOmit := utils.StructToMapOmit(req.DictTypeBase, nil, []string{"CreateTime"}, true)
	q := l.svcCtx.Dal.Query
	dictTypeBase, err := NewInfoLogic(l.ctx, l.svcCtx).Info(&types.IdReq{Id: req.DictID})
	if err != nil {
		return err
	}
	if sysDictType, err := q.SysDictType.WithContext(l.ctx).
		Where(q.SysDictType.DictID.Neq(req.DictID)).
		Where(q.SysDictType.DictType.Eq(req.DictType)).
		First(); err == nil && sysDictType != nil {
		return errx.BizErr("字典类型不能重复")
	}
	if _, err := q.SysDictType.WithContext(l.ctx).Where(l.svcCtx.Dal.Query.SysDictType.DictID.Eq(req.DictID)).Updates(toMapOmit); err != nil {
		return errx.GORMErr(err)
	}
	//如果修改字典类型名称 需要同步修改字典数据绑定字典
	if dictTypeBase.DictType != req.DictType {
		if _, err := q.SysDictDatum.WithContext(l.ctx).Where(q.SysDictDatum.DictType.Eq(dictTypeBase.DictType)).UpdateSimple(q.SysDictDatum.DictType.Value(req.DictType)); err != nil {
			return errx.GORMErr(err)
		}
	}
	return nil
}
