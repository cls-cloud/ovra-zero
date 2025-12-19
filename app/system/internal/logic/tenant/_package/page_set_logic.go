package _package

import (
	"context"
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

func (l *PageSetLogic) PageSet(req *types.PageSetTenantPackageReq) (resp *types.PageSetTenantPackageResp, err error) {
	offset := (req.PageNum - 1) * req.PageSize
	q := l.svcCtx.Dal.Query
	do := q.SysTenantPackage.WithContext(l.ctx)
	if req.PackageName != "" {
		do = do.Where(q.SysTenantPackage.PackageName.Like(fmt.Sprintf("%%%s%%", req.PackageName)))
	}
	result, count, err := do.Order(q.SysTenantPackage.CreateTime.Desc()).FindByPage(int(offset), int(req.PageSize))
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	list := make([]*types.TenantPackageBase, len(result))
	for i, item := range result {
		tenantPackage := new(types.TenantPackageBase)
		if err = copier.Copy(&tenantPackage, item); err != nil {
			return nil, err
		}
		list[i] = tenantPackage
	}
	resp = new(types.PageSetTenantPackageResp)
	resp.Rows = list
	resp.Total = count
	return
}
