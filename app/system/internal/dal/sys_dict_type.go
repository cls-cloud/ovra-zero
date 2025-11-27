package dal

import (
	"context"
	"system/internal/dal/model"
	"system/internal/dal/query"
	"system/internal/types"
	"toolkit/errx"
	"toolkit/utils"

	"gorm.io/gorm"
)

type SysDictTypeDal struct {
	query *query.Query
	db    *gorm.DB
}

func NewSysDictTypeDal(db *gorm.DB, query *query.Query) *SysDictTypeDal {
	return &SysDictTypeDal{
		db:    db,
		query: query,
	}
}

func (l *SysDictTypeDal) Insert(ctx context.Context, param *model.SysDictType) (err error) {
	su := l.query.SysDictType
	err = su.WithContext(ctx).Create(param)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysDictTypeDal) Update(ctx context.Context, param *model.SysDictType) (err error) {
	su := l.query.SysDictType
	if param.DictID == "" {
		return errx.BizErr("dictID is empty")
	}
	omit := utils.StructToMapOmit(param, nil, nil, true)
	_, err = su.WithContext(ctx).Where(su.DictID.Eq(param.DictID)).Updates(omit)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysDictTypeDal) Delete(ctx context.Context, id string) (err error) {
	su := l.query.SysDictType
	_, err = su.WithContext(ctx).Where(su.DictID.Eq(id)).Delete()
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysDictTypeDal) DeleteBatch(ctx context.Context, ids []string) (err error) {
	su := l.query.SysDictType
	_, err = su.WithContext(ctx).Where(su.DictID.In(ids...)).Delete()
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysDictTypeDal) SelectById(ctx context.Context, id string) (info *model.SysDictType, err error) {
	info = new(model.SysDictType)
	su := l.query.SysDictType
	data, err := su.WithContext(ctx).Where(su.DictID.Eq(id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	info = data
	return
}

func (l *SysDictTypeDal) PageSet(ctx context.Context, pageNum, pageSize int, query *types.DictTypeQuery) (total int64, list []*model.SysDictType, err error) {
	list = make([]*model.SysDictType, 0)
	limit := pageSize
	offset := (pageNum - 1) * pageSize
	su := l.query.SysDictType
	do := su.WithContext(ctx)
	result, count, err := do.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, errx.GORMErr(err)
	}
	list = result
	total = count
	return
}
