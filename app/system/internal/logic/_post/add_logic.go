package _post

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

func (l *AddLogic) Add(req *types.ModifyPostReq) error {
	id := utils.GetID()
	post := &model.SysPost{
		PostID:       id,
		PostCode:     req.PostCode,
		PostName:     req.PostName,
		PostSort:     req.PostSort,
		DeptID:       req.DeptID,
		Status:       req.Status,
		PostCategory: req.PostCategory,
		Remark:       req.Remark,
	}
	if err := l.svcCtx.Dal.SysPostDal.Insert(l.ctx, post); err != nil {
		return err
	}
	return nil
}
