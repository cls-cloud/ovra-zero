package visual

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

func (l *GetPageSetLogic) GetPageSet(req *types.VisualPageSetReq) (resp *types.VisualPageSetResp, err error) {
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageNo <= 0 {
		req.PageNo = 1
	}

	resp = new(types.VisualPageSetResp)
	offset := (req.PageNo - 1) * req.PageSize
	vc := l.svcCtx.Dal.Query.ReportVisual
	do := vc.WithContext(l.ctx)
	if req.Category != "" {
		do = do.Where(vc.Category.Eq(req.Category))
	}
	total, err := do.Count()
	if err != nil {
		return nil, err
	}
	if req.Category != "" {
		do = do.Where(vc.Category.Eq(req.Category))
	}
	categories, err := do.
		Order(vc.CreateTime.Desc()).
		Limit(int(req.PageSize)).
		Offset(int(offset)).
		Find()
	if err != nil {
		return nil, err
	}
	var rows []types.Visual
	for _, v := range categories {
		rows = append(rows, types.Visual{
			Uuid:          v.UUID,
			Title:         v.Title,
			BackgroundUrl: v.BackgroundURL,
			Thumbnail:     v.BackgroundURL,
			Category:      v.Category,
			Password:      v.Password,
			Status:        int64(v.Status),
			IsDeleted:     int64(v.IsDeleted),
		})
	}
	resp.PageNo = req.PageNo
	resp.PageSize = req.PageSize
	resp.TotalCount = total
	resp.TotalPage = (total + req.PageSize - 1) / req.PageSize
	resp.Rows = rows
	return
}
