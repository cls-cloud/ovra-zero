package _package

import (
	"context"
	"ovra/app/system/internal/dal/model"
	"ovra/toolkit/errx"
	"ovra/toolkit/utils"
	"strings"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.ModifyTenantPackageReq) error {
	q := l.svcCtx.Dal.Query
	oTenantPackage, err := q.SysTenantPackage.WithContext(l.ctx).
		Where(q.SysTenantPackage.PackageID.Neq(req.PackageID)).
		Where(q.SysTenantPackage.PackageName.Eq(req.PackageName)).
		First()
	if err == nil && oTenantPackage != nil {
		return errx.BizErr("租户套餐名称已存在")
	}
	menuIds := strings.Join(req.MenuIds, ",")
	tenantPackage := &model.SysTenantPackage{
		PackageID:         req.PackageID,
		PackageName:       req.PackageName,
		MenuIds:           menuIds,
		Remark:            req.Remark,
		MenuCheckStrictly: req.MenuCheckStrictly,
	}
	toMap := utils.StructToMapOmit(tenantPackage, nil, nil, true)
	if _, err := q.SysTenantPackage.WithContext(l.ctx).
		Where(q.SysTenantPackage.PackageID.Eq(req.PackageID)).
		Updates(toMap); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
