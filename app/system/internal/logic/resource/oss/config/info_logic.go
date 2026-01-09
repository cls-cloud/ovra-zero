package config

import (
	"context"
	"ovra/toolkit/errx"
	"strings"

	"github.com/jinzhu/copier"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoLogic) Info(req *types.IdReq) (resp *types.OssConfigBase, err error) {
	resp = new(types.OssConfigBase)
	q := l.svcCtx.Dal.Query
	ossConfig, err := q.SysOssConfig.WithContext(l.ctx).Where(q.SysOssConfig.OssConfigID.Eq(req.Id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	if err := copier.Copy(&resp, ossConfig); err != nil {
		return nil, err
	}
	resp.SecretKey = MaskKey(ossConfig.SecretKey)
	resp.AccessKey = MaskKey(ossConfig.AccessKey)
	return
}

func MaskKey(key string) string {
	if len(key) <= 5 {
		return key
	}
	return key[:5] + strings.Repeat("*", len(key)-5)
}

func IsMasked(key string) bool {
	if len(key) <= 5 {
		return false
	}
	suffix := key[5:]
	return strings.Trim(suffix, "*") == ""
}
