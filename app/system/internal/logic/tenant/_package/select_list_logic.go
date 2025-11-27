package _package

import (
	"context"
	"github.com/jinzhu/copier"
	"toolkit/errx"

	"system/internal/svc"
	"system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelectListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSelectListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectListLogic {
	return &SelectListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelectListLogic) SelectList() (resp []*types.TenantPackageBase, err error) {
	q := l.svcCtx.Dal.Query
	do := q.SysTenantPackage.WithContext(l.ctx)
	result, err := do.Order(q.SysTenantPackage.CreateTime.Desc()).Find()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = make([]*types.TenantPackageBase, len(result))
	for i, item := range result {
		tenantPackage := new(types.TenantPackageBase)
		if err = copier.Copy(&tenantPackage, item); err != nil {
			return nil, err
		}
		resp[i] = tenantPackage
	}
	return
}
