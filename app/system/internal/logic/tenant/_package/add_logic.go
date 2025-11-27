package _package

import (
	"context"
	"strings"
	"system/internal/dal/model"
	"toolkit/errx"
	"toolkit/utils"

	"system/internal/svc"
	"system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req *types.ModifyTenantPackageReq) error {
	menuIds := strings.Join(req.MenuIds, ",")
	tenantPackage := &model.SysTenantPackage{
		PackageID:         utils.GetID(),
		PackageName:       req.PackageName,
		MenuIds:           menuIds,
		Remark:            req.Remark,
		MenuCheckStrictly: req.MenuCheckStrictly,
	}
	if err := l.svcCtx.Dal.Query.SysTenantPackage.WithContext(l.ctx).Create(tenantPackage); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
