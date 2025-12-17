package config

import (
	"context"
	"resource/internal/dal/model"
	"toolkit/errx"
	"toolkit/utils"

	"resource/internal/svc"
	"resource/internal/types"

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

func (l *AddLogic) Add(req *types.ModifyOssConfigReq) error {
	ossConfig := &model.SysOssConfig{
		OssConfigID:  utils.GetID(),
		ConfigKey:    req.ConfigKey,
		AccessKey:    req.AccessKey,
		SecretKey:    req.SecretKey,
		BucketName:   req.BucketName,
		Prefix:       req.Prefix,
		Endpoint:     req.Endpoint,
		Domain:       req.Domain,
		IsHTTPS:      req.IsHTTPS,
		Region:       req.Region,
		AccessPolicy: req.AccessPolicy,
		Status:       req.Status,
		Ext1:         req.Ext1,
		Remark:       req.Remark,
	}
	if err := l.svcCtx.Dal.Query.SysOssConfig.WithContext(l.ctx).Create(ossConfig); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
