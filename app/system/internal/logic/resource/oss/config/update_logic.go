package config

import (
	"context"
	"encoding/json"
	"fmt"
	"ovra/toolkit/constants"
	"ovra/toolkit/errx"
	"ovra/toolkit/utils"

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

func (l *UpdateLogic) Update(req *types.ModifyOssConfigReq) error {
	q := l.svcCtx.Dal.Query
	if IsMasked(req.SecretKey) {
		req.SecretKey = ""
	}
	if IsMasked(req.AccessKey) {
		req.AccessKey = ""
	}
	toMapOmit := utils.StructToMapOmit(req.OssConfigBase, nil, []string{"CreateTime"}, true)
	if _, err := q.SysOssConfig.WithContext(l.ctx).Where(l.svcCtx.Dal.Query.SysOssConfig.OssConfigID.Eq(req.OssConfigID)).Updates(toMapOmit); err != nil {
		return errx.GORMErr(err)
	}
	config, err := q.SysOssConfig.WithContext(l.ctx).Where(q.SysOssConfig.OssConfigID.Eq(req.OssConfigID)).First()
	//更新redis
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = l.svcCtx.Rds.SetCtx(l.ctx, fmt.Sprintf(constants.OssConfigCache, config.ConfigKey), string(b))
	if err != nil {
		return err
	}
	if config.Status == "0" {
		err = l.svcCtx.Rds.SetCtx(l.ctx, constants.OssConfigDefaultCache, string(b))
		if err != nil {
			return err
		}
	}
	return nil
}
