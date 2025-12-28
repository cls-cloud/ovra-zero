package dal

import (
	"ovra/app/system/internal/config"
	"ovra/app/system/internal/dal/query"

	"gorm.io/gorm"
)

type Dal struct {
	Db              *gorm.DB
	Query           *query.Query
	Config          config.Config
	SysUserDal      *SysUserDal
	SysRoleDal      *SysRoleDal
	SysMenuDal      *SysMenuDal
	SysTenantDal    *SysTenantDal
	SysPostDal      *SysPostDal
	SysDeptDal      *SysDeptDal
	SysClientDal    *SysClientDal
	SysDictDatumDal *SysDictDatumDal
	SysDictTypeDal  *SysDictTypeDal
	SysNoticeDal    *SysNoticeDal
}

func NewDal(db *gorm.DB, query *query.Query, c config.Config) *Dal {
	return &Dal{
		Db:              db,
		Query:           query,
		Config:          c,
		SysUserDal:      NewSysUserDal(db, query),
		SysRoleDal:      NewSysRoleDal(db, query),
		SysMenuDal:      NewSysMenuDal(db, query),
		SysTenantDal:    NewSysTenantDal(db, query),
		SysPostDal:      NewSysPostDal(db, query),
		SysDeptDal:      NewSysDeptDal(db, query),
		SysClientDal:    NewSysClientDal(db, query),
		SysDictDatumDal: NewSysDictDatumDal(db, query),
		SysDictTypeDal:  NewSysDictTypeDal(db, query),
		SysNoticeDal:    NewSysNoticeDal(db, query),
	}
}
