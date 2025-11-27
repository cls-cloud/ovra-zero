package dal

import (
	"context"
	"system/internal/dal/model"
	"system/internal/dal/query"
	"system/internal/types"
	"toolkit/errx"
	"toolkit/utils"

	"gorm.io/gorm"
)

type SysPostDal struct {
	query *query.Query
	db    *gorm.DB
}

func NewSysPostDal(db *gorm.DB, query *query.Query) *SysPostDal {
	return &SysPostDal{
		db:    db,
		query: query,
	}
}

func (l *SysPostDal) Insert(ctx context.Context, param *model.SysPost) (err error) {
	su := l.query.SysPost
	err = su.WithContext(ctx).Create(param)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysPostDal) Update(ctx context.Context, param *model.SysPost) (err error) {
	su := l.query.SysPost
	if param.PostID == "" {
		return errx.BizErr("postID is empty")
	}
	omit := utils.StructToMapOmit(param, nil, nil, true)
	_, err = su.WithContext(ctx).Where(su.PostID.Eq(param.PostID)).Updates(omit)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysPostDal) Delete(ctx context.Context, id string) (err error) {
	err = l.query.Transaction(func(tx *query.Query) error {
		if _, err = tx.SysPost.WithContext(ctx).Where(tx.SysPost.PostID.Eq(id)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysUserPost.WithContext(ctx).Where(tx.SysUserPost.PostID.Eq(id)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		return nil
	})
	return
}

func (l *SysPostDal) DeleteBatch(ctx context.Context, ids []string) (err error) {

	err = l.query.Transaction(func(tx *query.Query) error {
		if _, err = tx.SysPost.WithContext(ctx).Where(tx.SysPost.PostID.In(ids...)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysUserPost.WithContext(ctx).Where(tx.SysUserPost.PostID.Eq(ids...)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		return nil
	})
	return
}

func (l *SysPostDal) SelectById(ctx context.Context, id string) (info *model.SysPost, err error) {
	info = new(model.SysPost)
	su := l.query.SysPost
	data, err := su.WithContext(ctx).Where(su.PostID.Eq(id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	info = data
	return
}

func (l *SysPostDal) PageSet(ctx context.Context, pageNum, pageSize int, query *types.PostQuery) (total int64, list []*model.SysPost, err error) {
	list = make([]*model.SysPost, 0)
	limit := pageSize
	offset := (pageNum - 1) * pageSize
	su := l.query.SysPost
	do := su.WithContext(ctx)
	result, count, err := do.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, errx.GORMErr(err)
	}
	list = result
	total = count
	return
}
