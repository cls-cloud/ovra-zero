package _post

import (
	"context"
	"toolkit/errx"
	"toolkit/utils"

	"system/internal/svc"
	"system/internal/types"

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
	toMapOmit := utils.StructToMapOmit(req.PostBase, nil, []string{"CreateTime"}, true)
	if _, err := l.svcCtx.Dal.Query.SysPost.WithContext(l.ctx).Where(l.svcCtx.Dal.Query.SysPost.PostID.Eq(req.PostID)).Updates(toMapOmit); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
