package _post

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

func (l *UpdateLogic) Update(req *types.ModifyPostReq) error {
	if err := l.svcCtx.Dal.SysPostDal.Update(l.ctx, &model.SysPost{
		PostID:       req.PostID,
		DeptID:       req.DeptID,
		PostCode:     req.PostCode,
		PostCategory: req.PostCategory,
		PostName:     req.PostName,
		PostSort:     req.PostSort,
		Status:       req.Status,
		Remark:       req.Remark,
	}); err != nil {
		return err
	}
	return nil
}
