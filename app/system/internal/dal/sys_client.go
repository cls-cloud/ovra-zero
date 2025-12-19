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

type SysClientDal struct {
	query *query.Query
	db    *gorm.DB
}

func NewSysClientDal(db *gorm.DB, query *query.Query) *SysClientDal {
	return &SysClientDal{
		db:    db,
		query: query,
	}
}

func (l *SysClientDal) Insert(ctx context.Context, param *model.SysClient) (err error) {
	su := l.query.SysClient
	err = su.WithContext(ctx).Create(param)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysClientDal) Update(ctx context.Context, param *model.SysClient) (err error) {
	su := l.query.SysClient
	if param.ClientID == "" {
		return errx.BizErr("clientID is empty")
	}
	omit := utils.StructToMapOmit(param, nil, nil, true)
	_, err = su.WithContext(ctx).Where(su.ClientID.Eq(param.ClientID)).Updates(omit)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysClientDal) Delete(ctx context.Context, id string) (err error) {
	su := l.query.SysClient
	_, err = su.WithContext(ctx).Where(su.ClientID.Eq(id)).Delete()
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysClientDal) DeleteBatch(ctx context.Context, ids []string) (err error) {
	su := l.query.SysClient
	_, err = su.WithContext(ctx).Where(su.ClientID.In(ids...)).Delete()
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysClientDal) SelectById(ctx context.Context, id string) (info *model.SysClient, err error) {
	info = new(model.SysClient)
	su := l.query.SysClient
	data, err := su.WithContext(ctx).Where(su.ClientID.Eq(id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	info = data
	return
}

func (l *SysClientDal) PageSet(ctx context.Context, pageNum, pageSize int, query *types.ClientQuery) (total int64, list []*model.SysClient, err error) {
	list = make([]*model.SysClient, 0)
	limit := pageSize
	offset := (pageNum - 1) * pageSize
	su := l.query.SysClient
	do := su.WithContext(ctx)
	result, count, err := do.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, errx.GORMErr(err)
	}
	list = result
	total = count
	return
}
