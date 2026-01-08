// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tree

import (
	"context"
	"ovra/app/demo/internal/dal/model"
	"ovra/app/demo/internal/svc"
	"ovra/app/demo/internal/types"

	"github.com/ovra-cloud/ovra-toolkit/errx"
	"github.com/ovra-cloud/ovra-toolkit/utils"

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

func (l *UpdateLogic) Update(req *types.ModifyTreeReq) error {
	if req.ID == 0 {
		return errx.BizErr("ID不能为空")
	}
	param := &model.TestTree{
		ID:       req.ID,
		ParentID: req.ParentID,
		DeptID:   utils.StrAtoi(req.DeptID),
		UserID:   utils.StrAtoi(req.UserID),
		TreeName: req.TreeName,
		Version:  int32(utils.StrAtoi(req.Version)),
	}
	if err := l.svcCtx.Dal.TestTree.Update(l.ctx, param); err != nil {
		return err
	}
	return nil
}
