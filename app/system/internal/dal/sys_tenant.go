package dal

import (
	"context"
	"ovra/app/system/internal/dal/model"
	"ovra/app/system/internal/dal/query"
	"ovra/app/system/internal/types"
	"ovra/toolkit/errx"
	"ovra/toolkit/utils"

	"gorm.io/gorm"
)

type SysTenantDal struct {
	query *query.Query
	db    *gorm.DB
}

func NewSysTenantDal(db *gorm.DB, query *query.Query) *SysTenantDal {
	return &SysTenantDal{
		db:    db,
		query: query,
	}
}

func (l *SysTenantDal) Insert(ctx context.Context, param *model.SysTenant) (err error) {
	su := l.query.SysTenant
	err = su.WithContext(ctx).Create(param)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysTenantDal) Update(ctx context.Context, param *model.SysTenant) (err error) {
	su := l.query.SysTenant
	if param.TenantID == "" {
		return errx.BizErr("tenantID is empty")
	}

	omit := utils.StructToMapOmit(param, nil, nil, true)
	_, err = su.WithContext(ctx).Where(su.TenantID.Eq(param.TenantID)).Updates(omit)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysTenantDal) UpdateStatus(ctx context.Context, id string, status string) (err error) {
	su := l.query.SysTenant
	if _, err = su.WithContext(ctx).Where(su.TenantID.Eq(id)).Update(su.Status, status); err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysTenantDal) Delete(ctx context.Context, id string) (err error) {
	su := l.query.SysTenant
	_, err = su.WithContext(ctx).Where(su.TenantID.Eq(id)).Delete()
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysTenantDal) DeleteBatch(ctx context.Context, ids []string) (err error) {
	su := l.query.SysTenant
	_, err = su.WithContext(ctx).Where(su.TenantID.In(ids...)).Delete()
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysTenantDal) SelectById(ctx context.Context, id string) (info *model.SysTenant, err error) {
	info = new(model.SysTenant)
	su := l.query.SysTenant
	data, err := su.WithContext(ctx).Where(su.TenantID.Eq(id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	info = data
	return
}

func (l *SysTenantDal) PageSet(ctx context.Context, pageNum, pageSize int, query *types.TenantQuery) (total int64, list []*model.SysTenant, err error) {
	list = make([]*model.SysTenant, 0)
	limit := pageSize
	offset := (pageNum - 1) * pageSize
	su := l.query.SysTenant
	do := su.WithContext(ctx)
	result, count, err := do.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, errx.GORMErr(err)
	}
	list = result
	total = count
	return
}
