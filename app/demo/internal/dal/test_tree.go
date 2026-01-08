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

type TestTreeDal struct {
	query *query.Query
	db    *gorm.DB
}

func NewTestTreeDal(db *gorm.DB, query *query.Query) *TestTreeDal {
	return &TestTreeDal{
		db:    db,
		query: query,
	}
}

func (l *TestTreeDal) Insert(ctx context.Context, param *model.TestTree) (err error) {
	su := l.query.TestTree
	err = su.WithContext(ctx).Create(param)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *TestTreeDal) Update(ctx context.Context, param *model.TestTree) (err error) {
	su := l.query.TestTree
	omit := utils.StructToMapOmit(param, []string{"Version"}, nil, true)
	_, err = su.WithContext(ctx).Where(su.ID.Eq(param.ID)).Updates(omit)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *TestTreeDal) Delete(ctx context.Context, id int64) (err error) {
	su := l.query.TestTree
	_, err = su.WithContext(ctx).Where(su.ID.Eq(id)).Delete()
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *TestTreeDal) DeleteBatch(ctx context.Context, ids []int64) (err error) {
	su := l.query.TestTree
	_, err = su.WithContext(ctx).Where(su.ID.In(ids...)).Delete()
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *TestTreeDal) SelectById(ctx context.Context, id int64) (*model.TestTree, error) {
	su := l.query.TestTree
	data, err := su.WithContext(ctx).Where(su.ID.Eq(id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	return data, nil
}
func (l *TestTreeDal) List(ctx context.Context, query *types.TreeQuery) (list []*model.TestTree, err error) {
	su := l.query.TestTree
	do := su.WithContext(ctx)
	if query.TreeName != "" {
		do = do.Where(su.TreeName.Like("%" + query.TreeName + "%"))
	}
	if query.UserID != "" {
		do = do.Where(su.UserID.Eq(utils.StrAtoi(query.UserID)))
	}
	if query.Version != "" {
		do = do.Where(su.Version.Eq(int32(utils.StrAtoi(query.Version))))
	}
	if query.DeptID != "" {
		do = do.Where(su.DeptID.Eq(utils.StrAtoi(query.DeptID)))
	}
	all, err := do.Order(su.CreateTime.Desc()).Find()
	if err != nil {
		return nil, errx.GORMErr(err)
	}

	// 如果没有 ParentID，直接返回
	if query.ParentID == "" {
		return all, nil
	}
	parentID := utils.StrAtoi(query.ParentID)
	group := make(map[int64][]*model.TestTree)
	for _, item := range all {
		group[item.ParentID] = append(group[item.ParentID], item)
	}
	var result []*model.TestTree
	queue := []int64{parentID}
	visited := make(map[int64]bool)
	for len(queue) > 0 {
		pid := queue[0]
		queue = queue[1:]

		if visited[pid] {
			continue
		}
		visited[pid] = true

		children := group[pid]
		for _, c := range children {
			result = append(result, c)
			queue = append(queue, c.ID)
		}
	}
	return result, nil
}
