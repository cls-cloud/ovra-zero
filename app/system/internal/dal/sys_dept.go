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

type SysDeptDal struct {
	query *query.Query
	db    *gorm.DB
}

func NewSysDeptDal(db *gorm.DB, query *query.Query) *SysDeptDal {
	return &SysDeptDal{
		db:    db,
		query: query,
	}
}

func (l *SysDeptDal) Insert(ctx context.Context, param *model.SysDept) (err error) {
	su := l.query.SysDept
	err = su.WithContext(ctx).Create(param)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysDeptDal) Update(ctx context.Context, param *model.SysDept) (err error) {
	su := l.query.SysDept
	if param.DeptID == "" {
		return errx.BizErr("deptID is empty")
	}
	omit := utils.StructToMapOmit(param, []string{"DeptCategory"}, nil, true)
	_, err = su.WithContext(ctx).Where(su.DeptID.Eq(param.DeptID)).Updates(omit)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysDeptDal) Delete(ctx context.Context, id string) (err error) {
	err = l.query.Transaction(func(tx *query.Query) error {
		if _, err = tx.SysDept.WithContext(ctx).Where(tx.SysDept.DeptID.Eq(id)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysRoleDept.WithContext(ctx).Where(tx.SysRoleDept.DeptID.Eq(id)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		return nil
	})
	return
}

func (l *SysDeptDal) DeleteBatch(ctx context.Context, ids []string) (err error) {
	err = l.query.Transaction(func(tx *query.Query) error {
		if _, err = tx.SysDept.WithContext(ctx).Where(tx.SysDept.DeptID.In(ids...)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysRoleDept.WithContext(ctx).Where(tx.SysRoleDept.DeptID.In(ids...)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		return nil
	})
	return
}

func (l *SysDeptDal) SelectById(ctx context.Context, id string) (*model.SysDept, error) {
	su := l.query.SysDept
	data, err := su.WithContext(ctx).Where(su.DeptID.Eq(id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	return data, nil
}

func (l *SysDeptDal) PageSet(ctx context.Context, pageNum, pageSize int, query *types.DeptQuery) (total int64, list []*model.SysDept, err error) {
	list = make([]*model.SysDept, 0)
	limit := pageSize
	offset := (pageNum - 1) * pageSize
	su := l.query.SysDept
	do := su.WithContext(ctx)
	result, count, err := do.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, errx.GORMErr(err)
	}
	list = result
	total = count
	return
}
