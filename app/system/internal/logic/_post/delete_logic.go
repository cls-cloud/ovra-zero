package _post

import (
	"context"
	"ovra/toolkit/errx"
	"strings"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

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
	q := l.svcCtx.Dal.Query
	if _, err := q.SysPost.WithContext(l.ctx).Where(q.SysPost.PostID.In(ids...)).Unscoped().Delete(); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
