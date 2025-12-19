package _package

import (
	"context"
	"ovra/toolkit/errx"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeStatusLogic {
	return &ChangeStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeStatusLogic) ChangeStatus(req *types.ChangeStatusTenantPackageReq) error {
	q := l.svcCtx.Dal.Query
	//查询正在使用该套餐的租户 如果有租户正在使用 需要先停用租户
	if req.Status == "1" {
		count, err := q.SysTenant.WithContext(l.ctx).Where(q.SysTenant.PackageID.Eq(req.PackageId), q.SysTenant.Status.Eq("0")).Count()
		if err != nil {
			return errx.GORMErr(err)
		}
		if count > 0 {
			return errx.BizErr("当前套餐正在被使用")
		}
	}
	_, err := q.SysTenantPackage.WithContext(l.ctx).Where(q.SysTenantPackage.PackageID.Eq(req.PackageId)).
		Update(q.SysTenantPackage.Status, req.Status)
	if err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
