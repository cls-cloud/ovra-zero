package sysrpclogic

import (
	"context"
	"fmt"
	"ovra/toolkit/errx"

	"ovra/app/system/internal/svc"
	"ovra/app/system/pb/system"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClientInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClientInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClientInfoLogic {
	return &ClientInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ClientInfoLogic) ClientInfo(req *system.ClientInfoReq) (*system.ClientInfoResp, error) {
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
	return &system.ClientInfoResp{
		Id:            client.ID,
		ClientId:      client.ClientID,
		ClientSecret:  client.ClientSecret,
		GrantType:     client.GrantType,
		ClientKey:     client.ClientKey,
		ActiveTimeout: client.ActiveTimeout,
		Timeout:       client.Timeout,
	}, nil
}
