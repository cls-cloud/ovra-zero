package dal

import (
	"system/internal/config"
	"system/internal/dal/query"

	"gorm.io/gorm"
)

type Dal struct {
	Db         *gorm.DB
	Query      *query.Query
	Config     config.Config
	SysUserDal *SysUserDal
	SysRoleDal *SysRoleDal
	SysMenuDal *SysMenuDal
}

func NewDal(db *gorm.DB, query *query.Query, c config.Config) *Dal {
	return &Dal{
		Db:         db,
		Query:      query,
		Config:     c,
		SysUserDal: NewSysUserDal(db, query),
		SysRoleDal: NewSysRoleDal(db, query),
		SysMenuDal: NewSysMenuDal(db, query),
	}
}
