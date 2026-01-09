package logininfor

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

func (l *PageSetLogic) PageSet(req *types.PageSetLogininforReq) (resp *types.PageSetLogininforResp, err error) {
	offset := (req.PageNum - 1) * req.PageSize
	q := l.svcCtx.Dal.Query
	do := q.SysLogininfor.WithContext(l.ctx)
	if req.Ipaddr != "" {
		do = do.Where(q.SysLogininfor.Ipaddr.Like(fmt.Sprintf("%%%s%%", req.Ipaddr)))
	}
	if req.UserName != "" {
		do = do.Where(q.SysLogininfor.UserName.Like(fmt.Sprintf("%%%s%%", req.UserName)))
	}
	if req.Status != "" {
		do = do.Where(q.SysLogininfor.Status.Eq(req.Status))
	}
	if req.Status != "" {
		do = do.Where(q.SysLogininfor.Status.Eq(req.Status))
	}
	if req.BeginTime != "" {
		beginTime, err := time.Parse(time.DateTime, req.BeginTime)
		if err != nil {
			return nil, errors.New("invalid beginTime format, expected: YYYY-MM-DD HH:mm:ss")
		}
		do = do.Where(q.SysLogininfor.LoginTime.Gte(beginTime))
	}
	if req.EndTime != "" {
		endTime, err := time.Parse(time.DateTime, req.EndTime)
		if err != nil {
			return nil, errors.New("invalid endTime format, expected: YYYY-MM-DD HH:mm:ss")
		}
		do = do.Where(q.SysLogininfor.LoginTime.Lte(endTime))
	}
	result, count, err := do.Order(q.SysLogininfor.LoginTime.Desc()).FindByPage(int(offset), int(req.PageSize))
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = new(types.PageSetLogininforResp)
	resp.Total = count
	list := make([]*types.LogininforBase, len(result))
	for i, item := range result {
		list[i] = new(types.LogininforBase)
		if err = copier.Copy(&list[i], item); err != nil {
			return nil, err
		}
		list[i].LoginTime = item.LoginTime.Format(time.DateTime)
		list[i].Ipaddr = DesensitizeIP(item.Ipaddr)
		list[i].LoginLocation = DesensitizeLocation(item.LoginLocation)
	}
	resp.Rows = list
	return
}

func DesensitizeIP(ipStr string) string {
	parts := strings.Split(ipStr, ".")
	if len(parts) == 4 {
		parts[3] = "*"
		parts[2] = "*"
		parts[1] = "*"
		return strings.Join(parts, ".")
	}
	return ipStr
}

func DesensitizeLocation(location string) string {
	parts := strings.Split(location, "|")
	if len(parts) == 4 {
		country := parts[0]
		province := parts[1]
		return fmt.Sprintf("%s|%s|**", country, province)
	}
	if len(parts) == 3 {
		country := parts[0]
		return fmt.Sprintf("%s|**", country)
	}
	return location
}
