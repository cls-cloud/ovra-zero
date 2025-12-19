package user

import (
	"context"
	"ovra/app/system/internal/dal/model"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStatusLogic {
	return &UpdateUserStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserStatusLogic) UpdateUserStatus(req *types.UpdateUserStatusReq) error {
	if err := l.svcCtx.Dal.SysUserDal.Update(l.ctx, &model.SysUser{
		UserID: req.UserId,
		Status: req.Status,
	}); err != nil {
		return err
	}
	return nil
}
