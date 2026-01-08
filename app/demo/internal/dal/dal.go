package dal

import (
	"ovra/app/demo/internal/config"
	"ovra/app/demo/internal/dal/query"

	"gorm.io/gorm"
)

type Dal struct {
	Db          *gorm.DB
	Query       *query.Query
	Config      config.Config
	TestDemoDal *TestDemoDal
	TestTree    *TestTreeDal
}

func NewDal(db *gorm.DB, query *query.Query, c config.Config) *Dal {
	return &Dal{
		Db:     db,
		Query:  query,
		Config: c,

		TestDemoDal: NewTestDemoDal(db, query),
		TestTree:    NewTestTreeDal(db, query),
	}
}
