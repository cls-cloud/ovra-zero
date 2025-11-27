package dal

import (
	"resource/internal/config"
	"resource/internal/dal/query"

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
