package visual

import (
	"context"
	"errors"
	"gorm.io/gorm"

	"report/internal/svc"
	"report/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.IdReq) (resp *types.Visual, err error) {
	reportVisual := l.svcCtx.Dal.Query.ReportVisual
	v, err := reportVisual.WithContext(l.ctx).Where(reportVisual.UUID.Eq(req.Id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("大屏不存在")
		}
		return nil, err
	}
	resp = &types.Visual{
		Uuid:          v.UUID,
		Title:         v.Title,
		BackgroundUrl: v.BackgroundURL,
		Thumbnail:     v.BackgroundURL,
		Category:      v.Category,
		Password:      v.Password,
		Status:        int64(v.Status),
		IsDeleted:     int64(v.IsDeleted),
	}
	return
}
