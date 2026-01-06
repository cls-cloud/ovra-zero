package client

import (
	"context"
	"ovra/toolkit/errx"
	"ovra/toolkit/utils"
	"strings"

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

func (l *UpdateLogic) Update(req *types.ModifyClientReq) error {
	if len(req.GrantTypeList) != 0 {
		req.GrantType = strings.Join(req.GrantTypeList, ",")
	}
	toMapOmit := utils.StructToMapOmit(req.ClientBase, nil, nil, true)
	if _, err := l.svcCtx.Dal.Query.SysClient.WithContext(l.ctx).Where(l.svcCtx.Dal.Query.SysClient.ID.Eq(req.ID)).Updates(toMapOmit); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
