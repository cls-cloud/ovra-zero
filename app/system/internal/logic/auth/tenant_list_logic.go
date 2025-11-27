package auth

import (
	"context"

	"system/internal/svc"
	"system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTenantListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantListLogic {
	return &TenantListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TenantListLogic) TenantList() (resp *types.TenantResp, err error) {
	resp = new(types.TenantResp)
	resp.TenantEnabled = l.svcCtx.Config.Tenant.Enabled
	if resp.TenantEnabled {
		resp.VoList = make([]types.TenantVo, 0)
		sysTenant := l.svcCtx.Dal.Query.SysTenant
		list, _ := sysTenant.WithContext(l.ctx).Find()
		for _, v := range list {
			vo := types.TenantVo{
				TenantId:    v.TenantID,
				CompanyName: v.CompanyName,
				Domain:      v.Domain,
			}
			resp.VoList = append(resp.VoList, vo)
		}
	}
	return
}
