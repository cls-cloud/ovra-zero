// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	"strings"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCascadeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCascadeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCascadeLogic {
	return &DeleteCascadeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCascadeLogic) DeleteCascade(req *types.IdsReq) error {
	ids := strings.Split(req.Ids, ",")
	dal := l.svcCtx.Dal
	if len(ids) > 0 {
		if err := dal.SysMenuDal.DeleteBatch(l.ctx, ids); err != nil {
			return err
		}
	}
	return nil
}
