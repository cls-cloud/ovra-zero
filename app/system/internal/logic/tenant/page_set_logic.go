package tenant

import (
	"context"
	"encoding/json"
	"fmt"
	"ovra/toolkit/errx"
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

type ListResp struct {
	Total int64                    `json:"total"`
	Rows  []map[string]interface{} `json:"rows"`
}

func (l *PageSetLogic) PageSet(req *types.PageSetTenantReq) (resp *ListResp, err error) {
	offset := (req.PageNum - 1) * req.PageSize
	q := l.svcCtx.Dal.Query
	do := q.SysTenant.WithContext(l.ctx)
	if req.TenantId != "" {
		do = do.Where(q.SysTenant.TenantID.Like(fmt.Sprintf("%%%s%%", req.TenantId)))
	}
	if req.CompanyName != "" {
		do = do.Where(q.SysTenant.CompanyName.Eq(req.CompanyName))
	}
	if req.ContactUserName != "" {
		do.Where(q.SysTenant.ContactUserName.Eq(fmt.Sprintf("%%%s%%", req.ContactUserName)))
	}
	result, count, err := do.Order(q.SysTenant.CreateTime.Desc()).FindByPage(int(offset), int(req.PageSize))
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = new(ListResp)
	resp.Total = count
	list := make([]map[string]interface{}, len(result))
	for i, item := range result {
		tenant := new(types.TenantBase)
		if err = copier.Copy(&tenant, item); err != nil {
			return nil, err
		}
		tenant.ExpireTime = item.ExpireTime.Format(time.DateTime)
		var toMap map[string]interface{}
		bs, err := json.Marshal(tenant)
		if err != nil {
			return nil, errx.BizErr("json 序列化失败")
		}
		err = json.Unmarshal(bs, &toMap)
		if err != nil {
			return nil, errx.BizErr("json 反序列化失败")
		}
		if item.ID == "1" {
			toMap["id"] = 1
		}
		list[i] = toMap
	}
	resp.Rows = list
	return
}
