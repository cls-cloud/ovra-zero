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

type SysDictDatumDal struct {
	query *query.Query
	db    *gorm.DB
}

func NewSysDictDatumDal(db *gorm.DB, query *query.Query) *SysDictDatumDal {
	return &SysDictDatumDal{
		db:    db,
		query: query,
	}
}

func (l *SysDictDatumDal) Insert(ctx context.Context, param *model.SysDictDatum) (err error) {
	su := l.query.SysDictDatum
	err = su.WithContext(ctx).Create(param)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysDictDatumDal) Update(ctx context.Context, param *model.SysDictDatum) (err error) {
	su := l.query.SysDictDatum
	if param.DictCode == "" {
		return errx.BizErr("dictCode is empty")
	}
	omit := utils.StructToMapOmit(param, nil, nil, true)
	_, err = su.WithContext(ctx).Where(su.DictCode.Eq(param.DictCode)).Updates(omit)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysDictDatumDal) Delete(ctx context.Context, id string) (err error) {
	su := l.query.SysDictDatum
	_, err = su.WithContext(ctx).Where(su.DictCode.Eq(id)).Delete()
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysDictDatumDal) DeleteBatch(ctx context.Context, ids []string) (err error) {
	su := l.query.SysDictDatum
	_, err = su.WithContext(ctx).Where(su.DictCode.In(ids...)).Delete()
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysDictDatumDal) SelectById(ctx context.Context, id string) (info *model.SysDictDatum, err error) {
	info = new(model.SysDictDatum)
	su := l.query.SysDictDatum
	data, err := su.WithContext(ctx).Where(su.DictCode.Eq(id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	info = data
	return
}

func (l *SysDictDatumDal) PageSet(ctx context.Context, pageNum, pageSize int, query *types.DictDataQuery) (total int64, list []*model.SysDictDatum, err error) {
	list = make([]*model.SysDictDatum, 0)
	limit := pageSize
	offset := (pageNum - 1) * pageSize
	su := l.query.SysDictDatum
	do := su.WithContext(ctx)
	result, count, err := do.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, errx.GORMErr(err)
	}
	list = result
	total = count
	return
}
