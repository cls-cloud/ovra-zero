package _type

import (
	"context"
	"ovra/app/system/internal/dal/model"
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

func (l *AddLogic) Add(req *types.ModifyDictTypeReq) error {
	dictType := &model.SysDictType{
		DictID:   utils.GetID(),
		DictName: req.DictName,
		DictType: req.DictType,
		Remark:   req.Remark,
	}
	if err := l.svcCtx.Dal.SysDictTypeDal.Insert(l.ctx, dictType); err != nil {
		return err
	}
	return nil
}
