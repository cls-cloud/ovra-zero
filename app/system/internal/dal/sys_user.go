package dal

import (
	"context"
	"ovra/app/system/internal/dal/model"
	"ovra/app/system/internal/dal/query"
	"ovra/app/system/internal/types"
	"ovra/toolkit/errx"
	"ovra/toolkit/helper"
	"ovra/toolkit/utils"

	"gorm.io/gorm"
)

type SysUserDal struct {
	query *query.Query
	db    *gorm.DB
}

func NewSysUserDal(db *gorm.DB, query *query.Query) *SysUserDal {
	return &SysUserDal{
		db:    db,
		query: query,
	}
}

func (l *SysUserDal) Insert(ctx context.Context, param *model.SysUser) (err error) {
	su := l.query.SysUser
	if param.Password != "" {
		param.Password = helper.Encrypt(param.Password)
	}
	err = su.WithContext(ctx).Create(param)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysUserDal) Update(ctx context.Context, param *model.SysUser) (err error) {
	su := l.query.SysUser
	if param.UserID == "" {
		return errx.BizErr("userID is empty")
	}
	if count, _ := su.WithContext(ctx).Where(su.UserID.Eq(param.UserID)).Count(); count != 1 {
		return errx.BizErr("用户不存在")
	}
	if param.Password != "" {
		param.Password = helper.Encrypt(param.Password)
	}
	omit := utils.StructToMapOmit(param, nil, nil, true)
	_, err = su.WithContext(ctx).Where(su.UserID.Eq(param.UserID)).Updates(omit)
	if err != nil {
		return errx.GORMErr(err)
	}
	return
}

func (l *SysUserDal) Delete(ctx context.Context, id string) (err error) {
	err = l.query.Transaction(func(tx *query.Query) error {
		if _, err = tx.SysUser.WithContext(ctx).Where(tx.SysUser.UserID.Eq(id)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysUserRole.WithContext(ctx).Where(tx.SysUserRole.UserID.Eq(id)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysUserPost.WithContext(ctx).Where(tx.SysUserPost.UserID.Eq(id)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		return nil
	})
	return
}

func (l *SysUserDal) DeleteBatch(ctx context.Context, ids []string) (err error) {
	err = l.query.Transaction(func(tx *query.Query) error {
		if _, err = tx.SysUser.WithContext(ctx).Where(tx.SysUser.UserID.In(ids...)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysUserRole.WithContext(ctx).Where(tx.SysUserRole.UserID.In(ids...)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		if _, err = tx.SysUserPost.WithContext(ctx).Where(tx.SysUserPost.UserID.In(ids...)).Delete(); err != nil {
			return errx.GORMErr(err)
		}
		return nil
	})
	return
}

func (l *SysUserDal) SelectById(ctx context.Context, id string) (info *model.SysUser, err error) {
	info = new(model.SysUser)
	su := l.query.SysUser
	data, err := su.WithContext(ctx).Where(su.UserID.Eq(id)).First()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	info = data
	return
}

func (l *SysUserDal) SelectByUserNameExit(ctx context.Context, userId, userName string) bool {
	su := l.query.SysUser
	do := su.WithContext(ctx)
	if userId != "" {
		do = do.Where(su.UserID.Neq(userId))
	}
	if count, _ := do.Where(su.UserName.Eq(userName)).Count(); count > 0 {
		return true
	}
	return false
}

func (l *SysUserDal) SelectByPhoneExit(ctx context.Context, userId, phone string) bool {
	su := l.query.SysUser
	do := su.WithContext(ctx)
	if userId != "" {
		do = do.Where(su.UserID.Neq(userId))
	}
	if count, _ := do.Where(su.Phonenumber.Eq(phone)).Count(); count > 0 {
		return true
	}
	return false
}

func (l *SysUserDal) SelectByEmailExit(ctx context.Context, userId, email string) bool {
	su := l.query.SysUser
	do := su.WithContext(ctx)
	if userId != "" {
		do = do.Where(su.UserID.Neq(userId))
	}
	if count, _ := do.Where(su.Email.Eq(email)).Count(); count > 0 {
		return true
	}
	return false
}

func (l *SysUserDal) SelectRoleIdsByUserId(ctx context.Context, userId string) (roleIds []string, err error) {
	roleIds = make([]string, 0)
	sur := l.query.SysUserRole
	err = sur.WithContext(ctx).Select().Select(sur.RoleID).Where(sur.UserID.Eq(userId)).Scan(&roleIds)
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	return
}
func (l *SysUserDal) SelectPostIdsByUserId(ctx context.Context, userId string) (postIds []string, err error) {
	postIds = make([]string, 0)
	sur := l.query.SysUserPost
	err = sur.WithContext(ctx).Select().Select(sur.PostID).Where(sur.UserID.Eq(userId)).Scan(&postIds)
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	return
}

func (l *SysUserDal) PageSet(ctx context.Context, pageNum, pageSize int, query *types.UserQuery) (total int64, list []*model.SysUser, err error) {
	list = make([]*model.SysUser, 0)
	limit := pageSize
	offset := (pageNum - 1) * pageSize
	su := l.query.SysUser
	do := su.WithContext(ctx)
	result, count, err := do.FindByPage(offset, limit)
	if err != nil {
		return 0, nil, errx.GORMErr(err)
	}
	list = result
	total = count
	return
}

func (l *SysUserDal) AddSysUserRoles(ctx context.Context, userId string, roleIds []string) (err error) {
	q := l.query
	err = q.Transaction(func(tx *query.Query) error {
		su := tx.SysUserRole
		if count, _ := su.WithContext(ctx).Where(su.UserID.Eq(userId)).Count(); count > 0 {
			_, err = su.WithContext(ctx).Where(su.UserID.Eq(userId)).Delete()
			if err != nil {
				return errx.GORMErr(err)
			}
		}
		rs := make([]*model.SysUserRole, len(roleIds))
		for i, roleId := range roleIds {
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

func (l *SysUserDal) AddSysUserPosts(ctx context.Context, userId string, postIds []string) (err error) {
	q := l.query
	err = q.Transaction(func(tx *query.Query) error {
		su := tx.SysUserPost
		if count, _ := su.WithContext(ctx).Where(su.UserID.Eq(userId)).Count(); count > 0 {
			_, err = su.WithContext(ctx).Where(su.UserID.Eq(userId)).Delete()
			if err != nil {
				return errx.GORMErr(err)
			}
		}
		rs := make([]*model.SysUserPost, len(postIds))
		for i, postId := range postIds {
			rs[i] = &model.SysUserPost{
				UserID: userId,
				PostID: postId,
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
