package client

import (
	"context"
	"encoding/json"
	"fmt"
	"ovra/toolkit/errx"

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

func (l *PageSetLogic) PageSet(req *types.PageSetClientReq) (resp *ListResp, err error) {
	offset := (req.PageNum - 1) * req.PageSize
	q := l.svcCtx.Dal.Query
	do := q.SysClient.WithContext(l.ctx)
	if req.ClientKey != "" {
		do = do.Where(q.SysClient.ClientKey.Like(fmt.Sprintf("%%%s%%", req.ClientKey)))
	}
	if req.ClientSecret != "" {
		do = do.Where(q.SysClient.ClientSecret.Like(fmt.Sprintf("%%%s%%", req.ClientSecret)))
	}
	if req.Status != "" {
		do = do.Where(q.SysClient.Status.Eq(req.Status))
	}
	result, count, err := do.Order(q.SysClient.ID.Asc()).FindByPage(int(offset), int(req.PageSize))
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = new(ListResp)
	resp.Total = count
	resp.Rows = make([]map[string]interface{}, len(result))
	resp.Total = count
	for i, item := range result {
		client := new(types.ClientBase)
		if err = copier.Copy(&client, item); err != nil {
			return nil, err
		}
		var toMap map[string]interface{}
		bs, err := json.Marshal(client)
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
		resp.Rows[i] = toMap
	}
	return
}
