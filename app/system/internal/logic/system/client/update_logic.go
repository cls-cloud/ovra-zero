package client

import (
	"context"
	"ovra/app/system/internal/dal/model"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.ModifyClientReq) error {
	if err := l.svcCtx.Dal.SysClientDal.Update(l.ctx, &model.SysClient{
		ID:            req.ID,
		ClientID:      req.ClientID,
		ClientKey:     req.ClientKey,
		ClientSecret:  req.ClientSecret,
		GrantType:     req.GrantType,
		DeviceType:    req.DeviceType,
		ActiveTimeout: req.ActiveTimeout,
		Timeout:       req.Timeout,
		Status:        req.Status,
	}); err != nil {
		return err
	}
	return nil
}
