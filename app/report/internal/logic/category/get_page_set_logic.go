package category

import (
	"context"
	"report/internal/svc"
	"report/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPageSetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPageSetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPageSetLogic {
	return &GetPageSetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *GetPageSetLogic) GetPageSet(req *types.PageReq) (resp *types.CategoryPageSetResp, err error) {
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageNo <= 0 {
		req.PageNo = 1
	}

	resp = new(types.CategoryPageSetResp)
	offset := (req.PageNo - 1) * req.PageSize
	vc := l.svcCtx.Dal.Query.ReportVisualCategory
	total, err := vc.WithContext(l.ctx).Count()
	if err != nil {
		return nil, err
	}
	categories, err := vc.WithContext(l.ctx).
		Order(vc.CategoryKey.Desc()).
		Offset(int(offset)).
		Limit(int(req.PageSize)).
		Find()
	if err != nil {
		return nil, err
	}
	var rows []types.VisualCategory
	for _, v := range categories {
		rows = append(rows, types.VisualCategory{
			CategoryKey:   v.CategoryKey,
			CategoryValue: v.CategoryValue,
			IsDeleted:     int64(v.IsDeleted),
			TenantUuid:    v.TenantUUID,
		})
	}
	resp.PageNo = req.PageNo
	resp.PageSize = req.PageSize
	resp.TotalCount = total
	resp.TotalPage = (total + req.PageSize - 1) / req.PageSize
	resp.Rows = rows
	return
}
