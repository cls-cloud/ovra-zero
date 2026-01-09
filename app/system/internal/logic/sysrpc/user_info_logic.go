package sysrpclogic

import (
	"context"
	"ovra/toolkit/errx"

	"ovra/app/system/internal/svc"
	"ovra/app/system/pb/system"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *system.UserInfoReq) (*system.UserInfoResp, error) {
	q := l.svcCtx.Query
	sysUser, err := q.SysUser.WithContext(l.ctx).Where(q.SysUser.UserName.Eq(in.Account)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}

	return &system.UserInfoResp{
		UserId:   sysUser.UserID,
		Username: sysUser.UserName,
		Password: sysUser.Password,
	}, nil
}
