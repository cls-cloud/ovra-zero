package config

import (
	"context"
	"errors"
	"fmt"
	"ovra/toolkit/errx"
	"time"

	"github.com/jinzhu/copier"

	"ovra/app/resource/internal/svc"
	"ovra/app/resource/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageSetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPageSetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageSetLogic {
	return &PageSetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageSetLogic) PageSet(req *types.PageSetOssConfigReq) (resp *types.PageSetOssConfigResp, err error) {
	offset := (req.PageNum - 1) * req.PageSize
	q := l.svcCtx.Dal.Query
	do := q.SysOssConfig.WithContext(l.ctx)
	if req.ConfigKey != "" {
		do = do.Where(q.SysOssConfig.ConfigKey.Like(fmt.Sprintf("%%%s%%", req.ConfigKey)))
	}
	if req.BucketName != "" {
		do = do.Where(q.SysOssConfig.BucketName.Like(fmt.Sprintf("%%%s%%", req.BucketName)))
	}
	if req.Status != "" {
		do = do.Where(q.SysOssConfig.Status.Eq(req.Status))
	}
	if req.BeginTime != "" {
		beginTime, err := time.Parse(time.DateTime, req.BeginTime)
		if err != nil {
			return nil, errors.New("invalid beginTime format, expected: YYYY-MM-DD HH:mm:ss")
		}
		do = do.Where(q.SysOssConfig.CreateTime.Gte(beginTime))
	}
	if req.EndTime != "" {
		endTime, err := time.Parse(time.DateTime, req.EndTime)
		if err != nil {
			return nil, errors.New("invalid endTime format, expected: YYYY-MM-DD HH:mm:ss")
		}
		do = do.Where(q.SysOssConfig.CreateTime.Lte(endTime))
	}
	result, count, err := do.Order(q.SysOssConfig.CreateTime.Desc()).FindByPage(int(offset), int(req.PageSize))
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = new(types.PageSetOssConfigResp)
	resp.Total = count
	list := make([]*types.OssConfigBase, len(result))
	for i, item := range result {
		list[i] = new(types.OssConfigBase)
		if err = copier.Copy(&list[i], item); err != nil {
			return nil, err
		}
		list[i].CreateTime = item.CreateTime.Format(time.DateTime)

		list[i].SecretKey = MaskKey(item.SecretKey)
		list[i].AccessKey = MaskKey(item.AccessKey)
	}
	resp.Rows = list
	return
}
