package dal

import (
	"ovra/app/resource/internal/config"
	"ovra/app/resource/internal/dal/query"

	"gorm.io/gorm"
)

type Dal struct {
	Db     *gorm.DB
	Query  *query.Query
	Config config.Config
}

func NewDal(db *gorm.DB, query *query.Query, c config.Config) *Dal {
	return &Dal{
		Db:     db,
		Query:  query,
		Config: c,
	}
}
