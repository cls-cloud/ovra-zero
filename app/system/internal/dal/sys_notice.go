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

type SysNoticeDal struct {
	query *query.Query
	db    *gorm.DB
}

func NewSysNoticeDal(db *gorm.DB, query *query.Query) *SysNoticeDal {
	return &SysNoticeDal{
		db:    db,
		query: query,
	}
}

func (l *SysNoticeDal) Insert(ctx context.Context, param *model.SysNotice) (err error) {
	su := l.query.SysNotice
	err = su.WithContext(ctx).Create(param)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysNoticeDal) Update(ctx context.Context, param *model.SysNotice) (err error) {
	su := l.query.SysNotice
	if param.NoticeID == "" {
		return errx.BizErr("noticeID is empty")
	}
	omit := utils.StructToMapOmit(param, nil, nil, true)
	_, err = su.WithContext(ctx).Where(su.NoticeID.Eq(param.NoticeID)).Updates(omit)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysNoticeDal) Delete(ctx context.Context, id string) (err error) {
	su := l.query.SysNotice
	_, err = su.WithContext(ctx).Where(su.NoticeID.Eq(id)).Delete()
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysNoticeDal) DeleteBatch(ctx context.Context, ids []string) (err error) {
	su := l.query.SysNotice
	_, err = su.WithContext(ctx).Where(su.NoticeID.In(ids...)).Delete()
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysNoticeDal) SelectById(ctx context.Context, id string) (info *model.SysNotice, err error) {
	info = new(model.SysNotice)
	su := l.query.SysNotice
	data, err := su.WithContext(ctx).Where(su.NoticeID.Eq(id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	info = data
	return
}

func (l *SysNoticeDal) PageSet(ctx context.Context, pageNum, pageSize int, query *types.NoticeQuery) (total int64, list []*model.SysNotice, err error) {
	list = make([]*model.SysNotice, 0)
	limit := pageSize
	offset := (pageNum - 1) * pageSize
	su := l.query.SysNotice
	do := su.WithContext(ctx)
	result, count, err := do.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, errx.GORMErr(err)
	}
	list = result
	total = count
	return
}
