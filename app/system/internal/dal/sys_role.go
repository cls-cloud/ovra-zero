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

type SysRoleDal struct {
	query *query.Query
	db    *gorm.DB
}

func NewSysRoleDal(db *gorm.DB, query *query.Query) *SysRoleDal {
	return &SysRoleDal{
		db:    db,
		query: query,
	}
}

func (l *SysRoleDal) Insert(ctx context.Context, param *model.SysRole) (err error) {
	su := l.query.SysRole
	err = su.WithContext(ctx).Create(param)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysRoleDal) Update(ctx context.Context, param *model.SysRole) (err error) {
	su := l.query.SysRole
	if param.RoleID == "" {
		return errx.BizErr("roleID is empty")
	}
	omit := utils.StructToMapOmit(param, nil, nil, true)
	_, err = su.WithContext(ctx).Where(su.RoleID.Eq(param.RoleID)).Updates(omit)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysRoleDal) UpdateStatus(ctx context.Context, id, status string) (err error) {
	su := l.query.SysRole
	_, err = su.WithContext(ctx).Where(su.RoleID.Eq(id)).Update(su.Status, status)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysRoleDal) Delete(ctx context.Context, id string) (err error) {
	if err = l.query.Transaction(func(tx *query.Query) error {
		if _, err = tx.SysRole.WithContext(ctx).Where(tx.SysRole.RoleID.Eq(id)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysRoleMenu.WithContext(ctx).Where(tx.SysRoleMenu.RoleID.Eq(id)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysRoleDept.WithContext(ctx).Where(tx.SysRoleDept.RoleID.Eq(id)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysUserRole.WithContext(ctx).Where(tx.SysUserRole.RoleID.Eq(id)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		return nil
	}); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}

func (l *SysRoleDal) DeleteBatch(ctx context.Context, ids []string) (err error) {
	if err = l.query.Transaction(func(tx *query.Query) error {
		if _, err = tx.SysRole.WithContext(ctx).Where(tx.SysRole.RoleID.In(ids...)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysRoleMenu.WithContext(ctx).Where(tx.SysRoleMenu.RoleID.In(ids...)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysRoleDept.WithContext(ctx).Where(tx.SysRoleDept.RoleID.In(ids...)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysUserRole.WithContext(ctx).Where(tx.SysUserRole.RoleID.In(ids...)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		return nil
	}); err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysRoleDal) SelectById(ctx context.Context, id string) (*model.SysRole, error) {
	su := l.query.SysRole
	data, err := su.WithContext(ctx).Where(su.RoleID.Eq(id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	return data, nil
}

func (l *SysRoleDal) SelectByRoleKeyExit(ctx context.Context, roleId, roleKey string) bool {
	su := l.query.SysRole
	do := su.WithContext(ctx)
	if roleId != "" {
		do = do.Where(su.RoleID.Neq(roleId))
	}
	if count, _ := do.Where(su.RoleKey.Eq(roleKey)).Count(); count > 0 {
		return true
	}
	return false
}

func (l *SysRoleDal) PageSet(ctx context.Context, pageNum, pageSize int, query *types.RoleQuery) (total int64, list []*model.SysRole, err error) {
	list = make([]*model.SysRole, 0)
	limit := pageSize
	offset := (pageNum - 1) * pageSize
	su := l.query.SysRole
	do := su.WithContext(ctx)
	result, count, err := do.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, errx.GORMErr(err)
	}
	list = result
	total = count
	return
}

func (l *SysRoleDal) AddSysRoleMenus(ctx context.Context, roleId string, menuIds []string) (err error) {
	q := l.query
	err = q.Transaction(func(tx *query.Query) error {
		su := tx.SysRoleMenu
		if count, _ := su.WithContext(ctx).Where(su.RoleID.Eq(roleId)).Count(); count > 0 {
			_, err = su.WithContext(ctx).Where(su.RoleID.Eq(roleId)).Delete()
			if err != nil {
				return errx.GORMErr(err)
			}
		}
		rs := make([]*model.SysRoleMenu, len(menuIds))
		for i, menuId := range menuIds {
			rs[i] = &model.SysRoleMenu{
				RoleID: roleId,
				MenuID: menuId,
			}
		}
		err = su.WithContext(ctx).CreateInBatches(rs, 100)
		if err != nil {
			return errx.GORMErr(err)
		}
		return nil
	})
	return
}

func (l *SysRoleDal) AddSysRoleUsers(ctx context.Context, roleId string, userIds []string) (err error) {
	q := l.query
	err = q.Transaction(func(tx *query.Query) error {
		su := tx.SysUserRole
		if count, _ := su.WithContext(ctx).Where(su.RoleID.Eq(roleId)).Count(); count > 0 {
			_, err = su.WithContext(ctx).Where(su.RoleID.Eq(roleId)).Delete()
			if err != nil {
				return errx.GORMErr(err)
			}
		}
		rs := make([]*model.SysUserRole, len(userIds))
		for i, userId := range userIds {
			rs[i] = &model.SysUserRole{
				UserID: userId,
				RoleID: roleId,
			}
		}
		err = su.WithContext(ctx).CreateInBatches(rs, 100)
		if err != nil {
			return errx.GORMErr(err)
		}
		return nil
	})
	return
}
