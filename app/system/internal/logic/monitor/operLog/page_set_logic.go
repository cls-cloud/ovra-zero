package operLog

import (
	"context"
	"errors"
	"fmt"
	"ovra/toolkit/errx"
	"strings"
	"time"

	"github.com/jinzhu/copier"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

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

func (l *PageSetLogic) PageSet(req *types.PageSetOperLogReq) (resp *types.PageSetOperLogResp, err error) {
	offset := (req.PageNum - 1) * req.PageSize
	q := l.svcCtx.Dal.Query
	do := q.SysOperLog.WithContext(l.ctx)
	if req.Title != "" {
		do = do.Where(q.SysOperLog.Title.Like(fmt.Sprintf("%%%s%%", req.Title)))
	}
	if req.OperName != "" {
		do = do.Where(q.SysOperLog.OperName.Like(fmt.Sprintf("%%%s%%", req.OperName)))
	}
	if req.BusinessType != 0 {
		do = do.Where(q.SysOperLog.BusinessType.Eq(req.BusinessType))
	}
	if req.OperIp != "" {
		do = do.Where(q.SysOperLog.OperIP.Eq(req.OperIp))
	}
	if req.BeginTime != "" {
		beginTime, err := time.Parse(time.DateTime, req.BeginTime)
		if err != nil {
			return nil, errors.New("invalid beginTime format, expected: YYYY-MM-DD HH:mm:ss")
		}
		do = do.Where(q.SysOperLog.OperTime.Gte(beginTime))
	}
	if req.EndTime != "" {
		endTime, err := time.Parse(time.DateTime, req.EndTime)
		if err != nil {
			return nil, errors.New("invalid endTime format, expected: YYYY-MM-DD HH:mm:ss")
		}
		do = do.Where(q.SysOperLog.OperTime.Lte(endTime))
	}
	result, count, err := do.Order(q.SysOperLog.OperTime.Asc()).FindByPage(int(offset), int(req.PageSize))
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = new(types.PageSetOperLogResp)
	resp.Total = count
	list := make([]*types.OperLogBase, len(result))
	for i, item := range result {
		list[i] = new(types.OperLogBase)
		if err = copier.Copy(&list[i], item); err != nil {
			return nil, err
		}
		list[i].OperTime = item.OperTime.Format(time.DateTime)
		list[i].OperIp = desensitizeIP(item.OperIP)
		list[i].OperLocation = desensitizeLocation(item.OperLocation)
	}
	resp.Rows = list
	return
}

func desensitizeIP(ipStr string) string {
	parts := strings.Split(ipStr, ".")
	if len(parts) == 4 {
		parts[3] = "*"
		parts[2] = "*"
		parts[1] = "*"
		return strings.Join(parts, ".")
	}
	return ipStr
}

func desensitizeLocation(location string) string {
	parts := strings.Split(location, "|")
	if len(parts) == 3 {
		country := parts[0]
		return fmt.Sprintf("%s***", country)
	}
	return location
}
