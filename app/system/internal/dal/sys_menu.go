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

type SysMenuDal struct {
	query *query.Query
	db    *gorm.DB
}

func NewSysMenuDal(db *gorm.DB, query *query.Query) *SysMenuDal {
	return &SysMenuDal{
		db:    db,
		query: query,
	}
}

func (l *SysMenuDal) Insert(ctx context.Context, param *model.SysMenu) (err error) {
	su := l.query.SysMenu
	err = su.WithContext(ctx).Create(param)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysMenuDal) Update(ctx context.Context, param *model.SysMenu) (err error) {
	su := l.query.SysMenu
	if param.MenuID == "" {
		return errx.BizErr("menuID is empty")
	}
	omit := utils.StructToMapOmit(param, nil, nil, true)
	_, err = su.WithContext(ctx).Where(su.MenuID.Eq(param.MenuID)).Updates(omit)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysMenuDal) Delete(ctx context.Context, id string) (err error) {
	err = l.query.Transaction(func(tx *query.Query) error {
		_, err = tx.SysMenu.WithContext(ctx).Where(tx.SysMenu.MenuID.Eq(id)).Delete()
		if err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysRoleMenu.WithContext(ctx).Where(tx.SysRoleMenu.MenuID.Eq(id)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		return nil
	})
	return
}

func (l *SysMenuDal) DeleteBatch(ctx context.Context, ids []string) (err error) {
	err = l.query.Transaction(func(tx *query.Query) error {
		_, err = tx.SysMenu.WithContext(ctx).Where(tx.SysMenu.MenuID.In(ids...)).Delete()
		if err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysRoleMenu.WithContext(ctx).Where(tx.SysRoleMenu.MenuID.In(ids...)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		return nil
	})
	return
}

func (l *SysMenuDal) SelectById(ctx context.Context, id string) (info *model.SysMenu, err error) {
	info = new(model.SysMenu)
	su := l.query.SysMenu
	data, err := su.WithContext(ctx).Where(su.MenuID.Eq(id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	info = data
	return
}

func (l *SysMenuDal) PageSet(ctx context.Context, pageNum, pageSize int, query *types.MenuQuery) (total int64, list []*model.SysMenu, err error) {
	list = make([]*model.SysMenu, 0)
	limit := pageSize
	offset := (pageNum - 1) * pageSize
	su := l.query.SysMenu
	do := su.WithContext(ctx)
	result, count, err := do.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, errx.GORMErr(err)
	}
	list = result
	total = count
	return
}
