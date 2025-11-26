package dal

import (
	"context"
	"system/internal/dal/model"
	"system/internal/dal/query"
	"system/internal/types"
	"time"
	"toolkit/errx"
	"toolkit/helper"
	"toolkit/utils"

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

func (l *SysUserDal) Insert(ctx context.Context, param *types.UserBase) (err error) {
	su := l.query.SysUser
	user := &model.SysUser{
		UserID:      utils.GetID(),
		UserName:    param.UserName,
		Password:    helper.Encrypt(param.Password),
		NickName:    param.NickName,
		Phonenumber: param.PhoneNumber,
		Email:       param.Email,
		Status:      param.Status,
		DeptID:      param.DeptID,
		Avatar:      param.Avatar,
		LoginIP:     param.LoginIp,
		UserType:    "00",
		LoginDate:   time.Now(),
		Remark:      param.Remark,
		Sex:         param.Sex,
	}
	err = su.WithContext(ctx).Create(user)
	return
}

func (l *SysUserDal) Update(ctx context.Context, param *types.UserBase) (err error) {
	su := l.query.SysUser
	if param.UserID == "" {
		return errx.BizErr("userId is empty")
	}
	omit := utils.StructToMapOmit(param, nil, nil, true)
	_, err = su.WithContext(ctx).Where(su.UserID.Eq(param.UserID)).Updates(omit)
	return err
}

func (l *SysUserDal) Delete() (list []*model.SysUser, err error) {

}

func (l *SysUserDal) DeleteBatch() (err error) {

}

func (l *SysUserDal) SelectById() (info *model.SysUser, err error) {

}

func (l *SysUserDal) PageSet() (total int64, list []*model.SysUser, err error) {

}
