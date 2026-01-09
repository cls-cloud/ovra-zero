package sysrpclogic

import (
	"context"

	"ovra/app/system/internal/svc"
	"ovra/app/system/pb/system"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTenantListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantListLogic {
	return &TenantListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TenantListLogic) TenantList(in *system.EmptyReq) (*system.TenantListResp, error) {
	resp := new(system.TenantListResp)
	resp.TenantEnable = l.svcCtx.Config.Tenant.Enabled
	if resp.TenantEnable {
		resp.List = make([]*system.TenantInfo, 0)
		sysTenant := l.svcCtx.Query.SysTenant
		list, _ := sysTenant.WithContext(l.ctx).Order(sysTenant.TenantID.Asc()).Find()
		for _, v := range list {
			vo := &system.TenantInfo{
				TenantId:    v.TenantID,
				CompanyName: v.CompanyName,
				Domain:      v.Domain,
			}
			resp.List = append(resp.List, vo)
		}
	}
	return resp, nil
}
