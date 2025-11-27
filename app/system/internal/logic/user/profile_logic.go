package user

import (
	"context"
	"github.com/jinzhu/copier"
	"strings"
	"system/internal/dal/model"
	"system/internal/svc"
	"system/internal/types"
	"time"
	"toolkit/auth"
	"toolkit/errx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProfileLogic) Profile() (resp *types.UserProfileResp, err error) {
	userId := auth.GetUserId(l.ctx)
	resp = new(types.UserProfileResp)
	q := l.svcCtx.Dal.Query
	var result struct {
		model.SysUser
		DeptName string `gorm:"column:dept_name"`
	}
	err = q.SysUser.WithContext(l.ctx).
		Select(
			q.SysUser.ALL,
			q.SysDept.DeptName, // 添加 dept_name 字段
		).
		LeftJoin(q.SysDept, q.SysDept.DeptID.EqCol(q.SysUser.DeptID)).
		Where(q.SysUser.UserID.Eq(userId)).
		Limit(1).
		Scan(&result)
	if err != nil {
		return nil, errx.GORMErrMsg(err, "用户不存在")
	}
	if err = copier.Copy(&resp.User, result); err != nil {
		return nil, err
	}
	resp.User.LoginDate = result.LoginDate.Format(time.DateTime)
	resp.User.DeptName = result.DeptName

	sysRole, err := q.SysRole.WithContext(l.ctx).
		LeftJoin(q.SysUserRole, q.SysUserRole.RoleID.EqCol(q.SysRole.RoleID)).
		Where(q.SysUserRole.UserID.Eq(userId)).
		Find()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	if err = copier.Copy(&resp.User.Roles, sysRole); err != nil {
		return nil, err
	}
	var roleNames []string
	for _, r := range sysRole {
		roleNames = append(roleNames, r.RoleName)
	}
	resp.RoleGroup = strings.Join(roleNames, ",")
	var postNames []string
	err = q.SysPost.WithContext(l.ctx).
		Select(q.SysPost.PostName).
		LeftJoin(q.SysUserPost, q.SysUserPost.PostID.EqCol(q.SysPost.PostID)).
		Where(q.SysUserPost.UserID.Eq(userId)).
		Scan(&postNames)
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp.PostGroup = strings.Join(postNames, ",")
	return
}
