package dept

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

func (l *UpdateLogic) Update(req *types.ModifyDeptReq) error {
	if err := l.svcCtx.Dal.SysDeptDal.Update(l.ctx, &model.SysDept{
		DeptID:       req.DeptID,
		ParentID:     req.ParentID,
		Ancestors:    req.Ancestors,
		DeptName:     req.DeptName,
		DeptCategory: req.DeptCategory,
		OrderNum:     req.OrderNum,
		Leader:       *req.Leader,
		Phone:        req.Phone,
		Email:        req.Email,
		Status:       req.Status,
	}); err != nil {
		return err
	}
	return nil
}
