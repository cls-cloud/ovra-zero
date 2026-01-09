package sysrpclogic

import (
	"context"
	"ovra/app/system/internal/dal/model"
	"ovra/toolkit/utils"
	"time"

	"ovra/app/system/internal/svc"
	"ovra/app/system/pb/system"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginInfoSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginInfoSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginInfoSaveLogic {
	return &LoginInfoSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginInfoSaveLogic) LoginInfoSave(in *system.LoginInfoReq) (*system.EmptyResp, error) {
	id := utils.GetID()
	if in.InfoId != "" {
		id = in.InfoId
	}
	logininfor := &model.SysLogininfor{
		InfoID:        id,
		TenantID:      in.TenantId,
		UserName:      in.Username,
		ClientKey:     in.ClientKey,
		DeviceType:    in.DeviceType,
		Ipaddr:        in.Ipaddr,
		LoginLocation: in.LoginLocation,
		Browser:       in.Browser,
		Os:            in.Os,
		Msg:           in.Msg,
		LoginTime:     time.Now(),
		Status:        in.Status,
	}
	err := l.svcCtx.Dal.Query.SysLogininfor.WithContext(l.ctx).Create(logininfor)
	return &system.EmptyResp{}, err
}
