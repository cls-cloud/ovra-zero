// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"fmt"
	"net/http"
	"ovra/app/auth/internal/types"
	"ovra/app/system/pb/system"
	"ovra/toolkit/auth"
	"ovra/toolkit/errx"
	"ovra/toolkit/ip"
	"ovra/toolkit/utils"

	"google.golang.org/api/idtoken"

	"ovra/app/auth/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
	resp = new(types.LoginResp)
	loginStatus := false
	loginInfoId := utils.GetID()
	if req.ClientId == "" {
		return nil, errx.AuthErr("clientId不能为空")
	}
	clientInfo, err := l.svcCtx.SysClient.ClientInfo(l.ctx, &system.ClientInfoReq{
		ClientId:  req.ClientId,
		GrantType: req.GrantType,
	})
	if err != nil {
		return nil, errx.AuthErr("Login failed, please contact the administrator")
	}

	if clientInfo.ClientKey == "pc" {
		resp, err = l.pcLogin(req, clientInfo, loginInfoId)
	}
	if err == nil {
		loginStatus = true
		l.LoginInfo("登录成功", loginStatus, req, loginInfoId, clientInfo)
	} else {
		l.LoginInfo(err.Error(), loginStatus, req, loginInfoId, clientInfo)
	}
	return
}

func (l *LoginLogic) pcLogin(req *types.LoginReq, clientInfo *system.ClientInfoResp, loginInfoId string) (*types.LoginResp, error) {
	loginInfo, err := l.svcCtx.SysClient.SysLogin(l.ctx, &system.LoginReq{
		Username:  req.Username,
		Password:  req.Password,
		TenantId:  req.TenantId,
		ClientId:  req.ClientId,
		GrantType: req.GrantType,
		Uuid:      req.Uuid,
		Code:      req.Code,
	})
	if err != nil {
		return nil, errx.AuthErr(err.Error())
	}
	var accessExpire = int64(clientInfo.Timeout)
	if int64(clientInfo.Timeout) == 0 {
		accessExpire = l.svcCtx.Config.JwtAuth.AccessExpire
	}
	accessSecret := l.svcCtx.Config.JwtAuth.AccessSecret

	ipStr, ua := ip.GetIPUa(l.r)
	name, version := ua.Browser()
	authMd5 := utils.AuthMd5(ipStr, name, version, ua.OS())

	userInfo := auth.UserInfo{
		UserId:   loginInfo.UserId,
		TenantId: loginInfo.TenantId,
		ClientId: loginInfo.ClientId,
		Timeout:  loginInfo.Timeout,
		DeptName: loginInfo.DeptName,
		UsMd5:    authMd5,
	}
	token, err := auth.GenerateToken(userInfo, accessSecret, accessExpire)
	if err != nil {
		return nil, err
	}
	key := ""
	if l.svcCtx.Config.JwtAuth.MultipleLoginDevices {
		key = fmt.Sprintf(auth.TokenKeyMd5, clientInfo.ClientId, loginInfo.UserId, authMd5)
	} else {
		key = fmt.Sprintf(auth.TokenKey, clientInfo.ClientId, loginInfo.UserId)
	}
	err = auth.NewAuth(l.svcCtx.Rds, &userInfo).SetToken(l.ctx, key, token, int64(clientInfo.ActiveTimeout), accessExpire, loginInfoId)
	if err != nil {
		return nil, err
	}
	return &types.LoginResp{AccessToken: token, ClientId: clientInfo.ClientId, ExpireIn: int64(clientInfo.Timeout), RefreshExpireIn: int64(clientInfo.ActiveTimeout)}, nil
}

func (l *LoginLogic) LoginInfo(msg string, status bool, req *types.LoginReq, loginInfoId string, clientInfo *system.ClientInfoResp) {
	location := "Unknown"
	cuRip, ua := ip.GetIPUa(l.r)
	ipAddr := "Unknown"
	logx.Info("请求 IP: ", cuRip)
	if !ip.IsPrivateIP(cuRip) {
		ipData, _ := ip.LookupIP(cuRip)
		location = ipData.Country
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
	if _, err := l.svcCtx.SysClient.LoginInfoSave(l.ctx, &system.LoginInfoReq{
		InfoId:        loginInfoId,
		Username:      req.Username,
		TenantId:      req.TenantId,
		ClientKey:     clientInfo.ClientKey,
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

// UserInfo 用来存储从 Google 返回的用户信息
type UserInfo struct {
	ID        string `json:"sub"`
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	FullName  string `json:"name"`
	Picture   string `json:"picture"`
}

// validateIDToken validates the id_token and returns user information
func validateIDToken(idToken, googleClientID string) (*UserInfo, error) {
	// Context for the token validation
	ctx := context.Background()

	// Validate the ID token and retrieve the payload
	payload, err := idtoken.Validate(ctx, idToken, googleClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate ID token: %v", err)
	}
	fmt.Println(payload)
	// Extract the claims map
	claims := payload.Claims

	// Extract relevant user info from the claims map
	userInfo := &UserInfo{
		ID:        claims["sub"].(string),         // The "sub" claim is the unique ID for the user
		Email:     claims["email"].(string),       // The "email" claim contains the user's email
		FirstName: claims["given_name"].(string),  // First name
		LastName:  claims["family_name"].(string), // Last name
		FullName:  claims["name"].(string),        // Full name
		Picture:   claims["picture"].(string),     // Profile picture URL
	}

	return userInfo, nil
}
