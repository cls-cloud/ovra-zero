package data

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

func (l *UpdateLogic) Update(req *types.ModifyDictDataReq) error {
	toMapOmit := utils.StructToMapOmit(req.DictDataBase, nil, []string{"CreateTime"}, true)
	if _, err := l.svcCtx.Dal.Query.SysDictDatum.WithContext(l.ctx).Where(l.svcCtx.Dal.Query.SysDictDatum.DictCode.Eq(req.DictCode)).Updates(toMapOmit); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
