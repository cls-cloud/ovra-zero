package client

import (
	"context"
	"ovra/app/system/internal/dal/model"
	"ovra/toolkit/errx"
	"ovra/toolkit/utils"
	"strings"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req *types.ModifyClientReq) error {
	if req.ClientID == "" {
		req.ClientID = strings.ReplaceAll(uuid.NewString(), "-", "")
	}
	client := &model.SysClient{
		ID:            utils.GetID(),
		ClientID:      req.ClientID,
		ClientKey:     req.ClientKey,
		ClientSecret:  req.ClientSecret,
		DeviceType:    req.DeviceType,
		GrantType:     req.GrantType,
		Status:        req.Status,
		Timeout:       req.Timeout,
		ActiveTimeout: req.ActiveTimeout,
	}
	if err := l.svcCtx.Dal.Query.SysClient.WithContext(l.ctx).Create(client); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
