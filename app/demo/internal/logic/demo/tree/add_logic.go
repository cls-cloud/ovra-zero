// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tree

import (
	"context"
	"ovra/app/demo/internal/dal/model"
	"ovra/app/demo/internal/svc"
	"ovra/app/demo/internal/types"

	"github.com/ovra-cloud/ovra-toolkit/utils"

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

func (l *AddLogic) Add(req *types.ModifyTreeReq) error {
	if err := l.svcCtx.Dal.TestTree.Insert(l.ctx, &model.TestTree{
		ParentID: req.ParentID,
		DeptID:   utils.StrAtoi(req.DeptID),
		UserID:   utils.StrAtoi(req.UserID),
		TreeName: req.TreeName,
		Version:  int32(utils.StrAtoi(req.Version)),
	}); err != nil {
		return err
	}
	return nil
}
