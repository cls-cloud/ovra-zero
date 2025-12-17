package logininforpclogic

import (
	"context"
	"monitor/internal/dal/model"
	"time"
	"toolkit/utils"

	"monitor/internal/svc"
	"monitor/pb/monitor"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveLogic {
	return &SaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SaveLogic) Save(in *monitor.LoginInfoReq) (*monitor.EmptyResp, error) {
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
	return &monitor.EmptyResp{}, err
}
