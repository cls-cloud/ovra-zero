package menu

import (
	"context"
	"ovra/toolkit/errx"
	"ovra/toolkit/utils"

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

func (l *UpdateLogic) Update(req *types.ModifyMenuReq) error {
	if req.MenuID == req.ParentID {
		return errx.BizErr("上级菜单不能为当前菜单")
	}
	toMapOmit := utils.StructToMapOmit(req.MenuBase, nil, []string{"CreateTime"}, true)
	if _, err := l.svcCtx.Dal.Query.SysMenu.WithContext(l.ctx).Where(l.svcCtx.Dal.Query.SysMenu.MenuID.Eq(req.MenuID)).Updates(toMapOmit); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
