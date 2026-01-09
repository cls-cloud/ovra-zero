package menu

import (
	"context"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.CodeReq) error {
	ids := strings.Split(req.Code, ",")
	dal := l.svcCtx.Dal
	if len(ids) > 0 {
		if err := dal.SysMenuDal.DeleteBatch(l.ctx, ids); err != nil {
			return err
		}
	}
	return nil
}
