package data

import (
	"context"
	"encoding/json"
	"ovra/app/system/internal/logic/system/dict/_type"
	"ovra/toolkit/constants"
	"ovra/toolkit/errx"

	"github.com/jinzhu/copier"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataByTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDataByTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataByTypeLogic {
	return &DataByTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DataByTypeLogic) DataByType(req *types.CodeReq) (resp []*types.DictDataBase, err error) {
	resp = make([]*types.DictDataBase, 0)
	cacheKey := constants.DictCache
	dataJson, err := l.svcCtx.Rds.HgetCtx(l.ctx, cacheKey, req.Code)
	if err == nil && len(dataJson) > 0 {
		type dictDataFromCache struct {
			DictCode   string `json:"dict_code"`
			TenantId   string `json:"tenant_id"`
			DictSort   int32  `json:"dict_sort"`
			DictLabel  string `json:"dict_label"`
			DictValue  string `json:"dict_value"`
			DictType   string `json:"dict_type"`
			CssClass   string `json:"css_class"`
			ListClass  string `json:"list_class"`
			IsDefault  string `json:"is_default"`
			Remark     string `json:"remark"`
			CreateTime string `json:"create_time"`
		}
		var cacheData []dictDataFromCache
		if err := json.Unmarshal([]byte(dataJson), &cacheData); err != nil {
			return nil, err
		}
		// 转换为目标结构体slice
		for _, d := range cacheData {
			resp = append(resp, &types.DictDataBase{
				DictCode:   d.DictCode,
				TenantId:   d.TenantId,
				DictSort:   d.DictSort,
				DictLabel:  d.DictLabel,
				DictValue:  d.DictValue,
				DictType:   d.DictType,
				CssClass:   d.CssClass,
				ListClass:  d.ListClass,
				IsDefault:  d.IsDefault,
				Remark:     d.Remark,
				CreateTime: d.CreateTime,
			})
		}
		return resp, nil
	}

	// 缓存未命中或解析失败，查数据库
	q := l.svcCtx.Dal.Query
	sysDictData, err := q.SysDictDatum.WithContext(l.ctx).Where(q.SysDictDatum.DictType.Eq(req.Code)).Find()
	if err != nil {
		return nil, errx.GORMErr(err)
	}

	// 拷贝到返回结构体
	err = copier.Copy(&resp, sysDictData)
	if err != nil {
		return nil, err
	}

	err = _type.NewRefreshCacheLogic(l.ctx, l.svcCtx).RefreshCache()
	return
}
