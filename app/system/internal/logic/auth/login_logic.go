package auth

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"monitor/pb/monitor"
	"net/http"
	"system/internal/svc"
	"system/internal/types"
	"time"
	"toolkit/auth"
	"toolkit/errx"
	"toolkit/helper"
	"toolkit/ip"
	"toolkit/utils"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	loginStatus := false
	loginInfoId := utils.GetID()

	resp, err = l.LoginHandler(req, loginInfoId, err)
	q := l.svcCtx.Dal.Query
	if err == nil {
		loginStatus = true
		l.LoginInfo("登陆成功", loginStatus, req, loginInfoId)
		//	更新登陆时间
		userMap := map[string]interface{}{
			"login_date": time.Now(),
			"login_ip":   ip.GetClientIP(l.r),
		}
		if _, err = q.SysUser.WithContext(l.ctx).Where(q.SysUser.UserName.Eq(req.Username), q.SysUser.TenantID.Eq(req.TenantId)).Updates(userMap); err != nil {
			logx.Error(err)
		}
	} else {
		l.LoginInfo(err.Error(), loginStatus, req, loginInfoId)
	}
	return
}

func (l *LoginLogic) LoginHandler(req *types.LoginReq, loginInfoId string, err error) (*types.LoginResp, error) {
	//判断clientId是否正常
	q := l.svcCtx.Dal.Query
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
	var accessExpire = int64(client.Timeout)
	if int64(client.Timeout) == 0 {
		accessExpire = l.svcCtx.Config.JwtAuth.AccessExpire
	}
	deptName := ""
	if user.DeptID != "" {
		dept, err := q.SysDept.WithContext(l.ctx).Where(q.SysDept.DeptID.Eq(user.DeptID)).First()
		if err != nil {
			return nil, errx.GORMErr(err)
		}
		deptName = dept.DeptName
	}

	accessSecret := l.svcCtx.Config.JwtAuth.AccessSecret
	userid := user.UserID

	ipStr, ua := ip.GetIPUa(l.r)
	name, version := ua.Browser()
	authMd5 := utils.AuthMd5(ipStr, name, version, ua.OS())

	userInfo := auth.UserInfo{
		UserId:   userid,
		TenantId: user.TenantID,
		ClientId: client.ClientID,
		Timeout:  client.Timeout,
		DeptName: deptName,
		UsMd5:    authMd5,
	}
	token, err := auth.GenerateToken(userInfo, accessSecret, accessExpire)
	if err != nil {
		return nil, err
	}
	key := ""
	if l.svcCtx.Config.JwtAuth.MultipleLoginDevices {
		key = fmt.Sprintf(auth.TokenKeyMd5, client.ClientID, userid, authMd5)
	} else {
		key = fmt.Sprintf(auth.TokenKey, client.ClientID, userid)
	}
	err = auth.NewAuth(l.svcCtx.Rds, &userInfo).SetToken(l.ctx, key, token, int64(client.ActiveTimeout), accessExpire, loginInfoId)
	if err != nil {
		return nil, err
	}
	return &types.LoginResp{AccessToken: token, ClientId: client.ClientID, ExpireIn: int64(client.Timeout), RefreshExpireIn: int64(client.ActiveTimeout)}, nil
}

func (l *LoginLogic) LoginInfo(msg string, status bool, req *types.LoginReq, loginInfoId string) {
	q := l.svcCtx.Dal.Query
	client, _ := q.SysClient.WithContext(l.ctx).Where(q.SysClient.ClientID.Eq(req.ClientId)).
		Where(q.SysClient.GrantType.Like(fmt.Sprintf("%%%s%%", req.GrantType))).
		First()
	location := "Unknown"
	cuRip, ua := ip.GetIPUa(l.r)
	ipAddr := "Unknown"
	logx.Info("请求 IP: ", cuRip)
	if !ip.IsPrivateIP(cuRip) {
		ipData, _ := ip.LookupIP(cuRip)
		location = fmt.Sprintf("%s|%s|%s|%s", ipData.Country, ipData.Region, ipData.City, ipData.ISP)
		ipAddr = cuRip
	} else {
		ipAddr = "内网IP"
	}
	browserName, _ := ua.Browser()
	deviceType := "PC"
	if ua.Mobile() {
		deviceType = "Mobile"
	}
	statusStr := "1"
	if status {
		statusStr = "0"
	}
	os := ip.ParseOS(ua.OS())
	if _, err := l.svcCtx.LoginInfoRpc.Save(l.ctx, &monitor.LoginInfoReq{
		InfoId:        loginInfoId,
		Username:      req.Username,
		TenantId:      req.TenantId,
		ClientKey:     client.ClientKey,
		DeviceType:    deviceType,
		Ipaddr:        ipAddr,
		LoginLocation: location,
		Browser:       browserName,
		Os:            os,
		Status:        statusStr,
		Msg:           msg,
	}); err != nil {
		logx.Error(err)
	}
}
