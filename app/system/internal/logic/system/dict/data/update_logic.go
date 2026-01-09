package data

import (
	"context"
	"ovra/app/system/internal/dal/model"
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

func (l *UpdateLogic) Update(req *types.ModifyDictDataReq) error {
	if err := l.svcCtx.Dal.SysDictDatumDal.Update(l.ctx, &model.SysDictDatum{
		DictCode:  req.DictCode,
		DictSort:  req.DictSort,
		DictLabel: req.DictLabel,
		DictValue: req.DictValue,
		DictType:  req.DictType,
		CSSClass:  req.CssClass,
		ListClass: req.ListClass,
		IsDefault: req.IsDefault,
		Remark:    req.Remark,
	}); err != nil {
		return err
	}
	return nil
}
