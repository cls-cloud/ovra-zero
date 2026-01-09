// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"ovra/app/system/pb/system"

	"ovra/app/auth/internal/svc"
	"ovra/app/auth/internal/types"

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
	listResp, err := l.svcCtx.SysClient.TenantList(l.ctx, &system.EmptyReq{})
	if err != nil {
		return
	}
	resp.TenantEnabled = listResp.TenantEnable
	if resp.TenantEnabled {
		resp.VoList = make([]types.TenantVo, 0)
		for _, v := range listResp.List {
			vo := types.TenantVo{
				TenantId:    v.TenantId,
				CompanyName: v.CompanyName,
				Domain:      v.Domain,
			}
			resp.VoList = append(resp.VoList, vo)
		}
	}
	return
}
