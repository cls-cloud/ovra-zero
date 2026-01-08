package dal

import (
	"context"
	"ovra/app/demo/internal/dal/model"
	"ovra/app/demo/internal/dal/query"
	"ovra/app/demo/internal/types"

	"github.com/ovra-cloud/ovra-toolkit/errx"
	"github.com/ovra-cloud/ovra-toolkit/utils"

	"gorm.io/gorm"
)

type TestDemoDal struct {
	query *query.Query
	db    *gorm.DB
}

func NewTestDemoDal(db *gorm.DB, query *query.Query) *TestDemoDal {
	return &TestDemoDal{
		db:    db,
		query: query,
	}
}

func (l *TestDemoDal) Insert(ctx context.Context, param *model.TestDemo) (err error) {
	su := l.query.TestDemo
	err = su.WithContext(ctx).Create(param)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *TestDemoDal) Update(ctx context.Context, param *model.TestDemo) (err error) {
	su := l.query.TestDemo
	omit := utils.StructToMapOmit(param, nil, nil, true)
	_, err = su.WithContext(ctx).Where(su.ID.Eq(param.ID)).Updates(omit)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *TestDemoDal) Delete(ctx context.Context, id int64) (err error) {
	su := l.query.TestDemo
	_, err = su.WithContext(ctx).Where(su.ID.Eq(id)).Delete()
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *TestDemoDal) DeleteBatch(ctx context.Context, ids []int64) (err error) {
	su := l.query.TestDemo
	_, err = su.WithContext(ctx).Where(su.ID.In(ids...)).Delete()
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *TestDemoDal) SelectById(ctx context.Context, id int64) (*model.TestDemo, error) {
	su := l.query.TestDemo
	data, err := su.WithContext(ctx).Where(su.ID.Eq(id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	return data, nil
}

func (l *TestDemoDal) PageSet(ctx context.Context, pageNum, pageSize int, query *types.DemoQuery) (total int64, list []*model.TestDemo, err error) {
	list = make([]*model.TestDemo, 0)
	limit := pageSize
	offset := (pageNum - 1) * pageSize
	su := l.query.TestDemo
	do := su.WithContext(ctx)
	if query.TestKey != "" {
		do = do.Where(su.TestKey.Like("%" + query.TestKey + "%"))
	}
	if query.OrderNum != "" {
		do = do.Where(su.OrderNum.Eq(int32(utils.StrAtoi(query.OrderNum))))
	}
	if query.Value != "" {
		do = do.Where(su.Value.Eq(query.Value))
	}
	if query.Version != "" {
		do = do.Where(su.Version.Eq(int32(utils.StrAtoi(query.Version))))
	}
	result, count, err := do.Order(su.OrderNum.Asc(), su.CreateTime.Desc()).FindByPage(offset, limit)
	if err != nil {
		return 0, nil, errx.GORMErr(err)
	}
	list = result
	total = count
	return
}
