// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tree

import (
	"context"
	"ovra/app/demo/internal/svc"
	"ovra/app/demo/internal/types"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.ListTreeReq) (resp []*types.TreeBase, err error) {
	list, err := l.svcCtx.Dal.TestTree.List(l.ctx, &req.TreeQuery)
	if err != nil {
		return nil, err
	}
	resp = make([]*types.TreeBase, len(list))
	for i, item := range list {
		resp[i] = &types.TreeBase{
			ID:       item.ID,
			ParentID: item.ParentID,
			DeptID:   strconv.FormatInt(item.DeptID, 10),
			UserID:   strconv.FormatInt(item.UserID, 10),
			TreeName: item.TreeName,
			Version:  strconv.Itoa(int(item.Version)),
		}
	}
	return
}
