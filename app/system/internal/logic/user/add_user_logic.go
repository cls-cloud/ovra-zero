package user

import (
	"context"
	"system/internal/dal/model"
	"time"
	"toolkit/errx"
	"toolkit/utils"

	"system/internal/svc"
	"system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserLogic {
	return &AddUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddUserLogic) AddUser(req *types.AddOrUpdateUserReq) error {
	userId := utils.GetID()
	dal := l.svcCtx.Dal
	if req.UserName != "" {
		if exit := dal.SysUserDal.SelectByUserNameExit(l.ctx, "", req.UserName); exit {
			return errx.BizErr("用户名已经存在")
		}
	}
	if req.PhoneNumber != "" {
		if exit := dal.SysUserDal.SelectByPhoneExit(l.ctx, "", req.PhoneNumber); exit {
			return errx.BizErr("手机号已经存在")
		}
	}
	if req.Email != "" {
		if exit := dal.SysUserDal.SelectByEmailExit(l.ctx, "", req.Email); exit {
			return errx.BizErr("邮箱已经存在")
		}
	}

	user, err := l.genUser(req, userId)
	if err != nil {
		return err
	}
	if err := dal.SysUserDal.Insert(l.ctx, user); err != nil {
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

func (l *AddUserLogic) genUser(req *types.AddOrUpdateUserReq, userId string) (*model.SysUser, error) {
	if req.UserID != "" {
		return nil, errx.BizErr("用户ID不为空")
	}
	if req.Password == "" {
		return nil, errx.BizErr("密码不能为空")
	}
	user := &model.SysUser{
		UserID:      userId,
		UserName:    req.UserName,
		Password:    req.Password,
		NickName:    req.NickName,
		Phonenumber: req.PhoneNumber,
		Email:       req.Email,
		Status:      req.Status,
		DeptID:      req.DeptID,
		Avatar:      req.Avatar,
		LoginIP:     req.LoginIp,
		UserType:    "00",
		LoginDate:   time.Now(),
		Remark:      req.Remark,
		Sex:         req.Sex,
	}
	if req.TenantID != "" {
		user.TenantID = req.TenantID
	}
	return user, nil
}
