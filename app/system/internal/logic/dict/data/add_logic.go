package data

import (
	"context"
	"ovra/app/system/internal/dal/model"
	"ovra/toolkit/errx"
	"ovra/toolkit/utils"

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

func (l *AddLogic) Add(req *types.ModifyDictDataReq) error {
	dictData := &model.SysDictDatum{
		DictCode:  utils.GetID(),
		DictLabel: req.DictLabel,
		DictSort:  req.DictSort,
		DictType:  req.DictType,
		DictValue: req.DictValue,
		IsDefault: req.IsDefault,
		ListClass: req.ListClass,
		CSSClass:  req.CssClass,
	}
	q := l.svcCtx.Dal.Query.SysDictDatum
	if err := q.WithContext(l.ctx).Create(dictData); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
