package user

import (
	"context"
	"ovra/app/system/internal/dal/model"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"ovra/toolkit/errx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.AddOrUpdateUserReq) error {
	userId := req.UserID
	dal := l.svcCtx.Dal
	if req.UserName != "" {
		if exit := dal.SysUserDal.SelectByUserNameExit(l.ctx, userId, req.UserName); exit {
			return errx.BizErr("用户名已经存在")
		}
	}
	if req.PhoneNumber != "" {
		if exit := dal.SysUserDal.SelectByPhoneExit(l.ctx, userId, req.PhoneNumber); exit {
			return errx.BizErr("手机号已经存在")
		}
	}
	if req.Email != "" {
		if exit := dal.SysUserDal.SelectByEmailExit(l.ctx, userId, req.Email); exit {
			return errx.BizErr("邮箱已经存在")
		}
	}
	if err := dal.SysUserDal.Update(l.ctx, &model.SysUser{
		UserID:      userId,
		DeptID:      req.DeptID,
		UserName:    req.UserName,
		NickName:    req.NickName,
		UserType:    req.UserType,
		Email:       req.Email,
		Phonenumber: req.PhoneNumber,
		Sex:         req.Sex,
		Avatar:      req.Avatar,
		Status:      req.Status,
		Remark:      req.Remark,
	}); err != nil {
		return err
	}
	if len(req.RoleIds) != 0 {
		if err := dal.SysUserDal.AddSysUserRoles(l.ctx, userId, req.RoleIds); err != nil {
			return err
		}
	}
	if len(req.PostIds) != 0 {
		if err := dal.SysUserDal.AddSysUserPosts(l.ctx, userId, req.PostIds); err != nil {
			return err
		}
	}
	return nil
}
