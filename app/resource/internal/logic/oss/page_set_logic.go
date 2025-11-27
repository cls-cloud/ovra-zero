package oss

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"time"
	"toolkit/errx"

	"resource/internal/svc"
	"resource/internal/types"

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

func (l *PageSetLogic) PageSet(req *types.PageSetOssReq) (resp *types.PageSetOssResp, err error) {
	offset := (req.PageNum - 1) * req.PageSize
	q := l.svcCtx.Dal.Query
	do := q.SysOss.WithContext(l.ctx)
	if req.FileName != "" {
		do = do.Where(q.SysOss.FileName.Like(fmt.Sprintf("%%%s%%", req.FileName)))
	}
	if req.OriginalName != "" {
		do = do.Where(q.SysOss.OriginalName.Like(fmt.Sprintf("%%%s%%", req.OriginalName)))
	}
	if req.FileSuffix != "" {
		do = do.Where(q.SysOss.FileSuffix.Eq(req.FileSuffix))
	}
	if req.Service != "" {
		do = do.Where(q.SysOss.Service.Eq(req.Service))
	}
	if req.BeginTime != "" {
		beginTime, err := time.Parse(time.DateTime, req.BeginTime)
		if err != nil {
			return nil, errors.New("invalid beginTime format, expected: YYYY-MM-DD HH:mm:ss")
		}
		do = do.Where(q.SysOss.CreateTime.Gte(beginTime))
	}
	if req.EndTime != "" {
		endTime, err := time.Parse(time.DateTime, req.EndTime)
		if err != nil {
			return nil, errors.New("invalid endTime format, expected: YYYY-MM-DD HH:mm:ss")
		}
		do = do.Where(q.SysOss.CreateTime.Lte(endTime))
	}
	result, count, err := do.Order(q.SysOss.CreateTime.Desc()).FindByPage(int(offset), int(req.PageSize))
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = new(types.PageSetOssResp)
	resp.Total = count
	list := make([]*types.OssBase, len(result))
	for i, item := range result {
		list[i] = new(types.OssBase)
		if err = copier.Copy(&list[i], item); err != nil {
			return nil, err
		}
		list[i].CreateTime = item.CreateTime.Format(time.DateTime)
	}
	resp.Rows = list
	return
}
