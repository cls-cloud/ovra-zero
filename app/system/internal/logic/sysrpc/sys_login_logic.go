package sysrpclogic

import (
	"context"
	"fmt"
	"ovra/toolkit/errx"
	"ovra/toolkit/helper"

	"ovra/app/system/internal/svc"
	"ovra/app/system/pb/system"

	"github.com/zeromicro/go-zero/core/logx"
)

type SysLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSysLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SysLoginLogic {
	return &SysLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SysLoginLogic) SysLogin(req *system.LoginReq) (*system.LoginResp, error) {
	//判断clientId是否正常
	q := l.svcCtx.Query
	if req.ClientId == "" {
		return nil, errx.AuthErr("clientId不能为空")
	}
	client, err := q.SysClient.WithContext(l.ctx).Where(q.SysClient.ClientID.Eq(req.ClientId), q.SysClient.Status.Eq("0")).
		Where(q.SysClient.GrantType.Like(fmt.Sprintf("%%%s%%", req.GrantType))).
		First()
	if err != nil {
		return nil, errx.GORMErrMsg(err, fmt.Sprintf("客户端id: %s 认证类型：%s 认证异常!", req.ClientId, req.GrantType))
	}
	_, err = q.SysTenant.WithContext(l.ctx).Where(q.SysTenant.TenantID.Eq(req.TenantId), q.SysTenant.Status.Eq("0")).First()
	if err != nil {
		return nil, errx.GORMErrMsg(err, "租户不存在或已停用")
	}
	if l.svcCtx.Config.Captcha.Enabled {
		data, err := l.svcCtx.Rds.GetCtx(l.ctx, "captcha:"+req.Uuid)
		if err != nil {
			return nil, errx.AuthErr("验证码失效")
		}
		if data != req.Code {
			//验证码错误则重新获取验证码
			_, _ = l.svcCtx.Rds.DelCtx(l.ctx, "captcha:"+req.Uuid)
			return nil, errx.AuthErr("验证码验证失败")
		}
	}
	sysUser := q.SysUser
	user, err := sysUser.WithContext(l.ctx).Where(sysUser.UserName.Eq(req.Username)).
		Where(sysUser.TenantID.Eq(req.TenantId)).
		First()

	if err != nil {
		return nil, errx.GORMErrMsg(err, "用户不存在")
	}
	//判断密码是否匹配
	if !helper.Verify(user.Password, req.Password) {
		return nil, errx.AuthErr("密码验证失败")
	}
	deptName := ""
	if user.DeptID != "" {
		dept, err := q.SysDept.WithContext(l.ctx).Where(q.SysDept.DeptID.Eq(user.DeptID)).First()
		if err != nil {
			return nil, errx.GORMErr(err)
		}
		deptName = dept.DeptName
	}
	return &system.LoginResp{
		UserId:   user.UserID,
		TenantId: user.TenantID,
		ClientId: client.ClientID,
		Timeout:  client.Timeout,
		DeptName: deptName,
	}, nil
}
