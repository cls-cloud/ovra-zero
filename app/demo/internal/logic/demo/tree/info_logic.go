// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tree

import (
	"context"
	"ovra/app/demo/internal/svc"
	"ovra/app/demo/internal/types"
	"strconv"

	"github.com/ovra-cloud/ovra-toolkit/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoLogic) Info(req *types.IdReq) (resp *types.TreeBase, err error) {
	data, err := l.svcCtx.Dal.TestTree.SelectById(l.ctx, utils.StrAtoi(req.Id))
	if err != nil {
		return nil, err
	}
	resp = &types.TreeBase{
		ID:       data.ID,
		ParentID: data.ParentID,
		DeptID:   strconv.FormatInt(data.DeptID, 10),
		UserID:   strconv.FormatInt(data.UserID, 10),
		TreeName: data.TreeName,
		Version:  strconv.Itoa(int(data.Version)),
	}
	return
}
