package user

import (
	"context"
	"system/internal/dal/model"
	"system/internal/svc"
	"system/internal/types"
	"toolkit/errx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPwdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPwdLogic {
	return &ResetPwdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPwdLogic) ResetPwd(req *types.ResetPwdReq) error {
	if len(req.Password) < 5 || len(req.Password) > 20 {
		return errx.BizErr("密码长度为5 - 20")
	}
	if err := l.svcCtx.Dal.SysUserDal.Update(l.ctx, &model.SysUser{
		UserID:   req.UserId,
		Password: req.Password,
	}); err != nil {
		return err
	}
	return nil
}
