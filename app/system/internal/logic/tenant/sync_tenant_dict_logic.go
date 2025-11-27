package tenant

import (
	"context"
	"toolkit/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"system/internal/svc"
)

type SyncTenantDictLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncTenantDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncTenantDictLogic {
	return &SyncTenantDictLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncTenantDictLogic) SyncTenantDict() error {
	baseTenantId := "000000"
	q := l.svcCtx.Dal.Query
	sysTenants, err := q.SysTenant.WithContext(l.ctx).Where(q.SysTenant.Status.Eq("0"), q.SysTenant.TenantID.Neq(baseTenantId)).Find()
	if err != nil {
		return err
	}
	for _, tenant := range sysTenants {
		sysDictTypes, err := q.SysDictType.WithContext(l.ctx).Where(q.SysDictType.TenantID.Eq(baseTenantId)).Find()
		if err != nil {
			return err
		}
		_, err = q.SysDictType.WithContext(l.ctx).Where(q.SysDictType.TenantID.Eq(tenant.TenantID)).Delete()
		if err != nil {
			return err
		}
		for i := range sysDictTypes {
			sysDictTypes[i].DictID = utils.GetID()
			sysDictTypes[i].TenantID = tenant.TenantID
		}
		err = q.SysDictType.WithContext(l.ctx).CreateInBatches(sysDictTypes, len(sysDictTypes))
		if err != nil {
			return err
		}

		sysDictDatums, err := q.SysDictDatum.WithContext(l.ctx).Where(q.SysDictDatum.TenantID.Eq(baseTenantId)).Find()
		if err != nil {
			return err
		}
		_, err = q.SysDictDatum.WithContext(l.ctx).Where(q.SysDictDatum.TenantID.Eq(tenant.TenantID)).Delete()
		if err != nil {
			return err
		}
		for i := range sysDictDatums {
			sysDictDatums[i].DictCode = utils.GetID()
			sysDictDatums[i].TenantID = tenant.TenantID
		}
		err = q.SysDictDatum.WithContext(l.ctx).CreateInBatches(sysDictDatums, len(sysDictDatums))
	}
	return nil
}
