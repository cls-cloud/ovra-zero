package user

import (
	"context"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"ovra/toolkit/errx"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type QueryUserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryUserDetailLogic {
	return &QueryUserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryUserDetailLogic) QueryUserDetail(req *types.IdReq) (resp *types.UserDetailResp, err error) {
	dal := l.svcCtx.Dal
	q := dal.Query
	resp = new(types.UserDetailResp)
	// 查询系统所有角色
	roles, err := q.SysRole.WithContext(l.ctx).Find()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	if err = copier.Copy(&resp.Roles, roles); err != nil {
		return nil, err
	}
	posts, err := q.SysPost.WithContext(l.ctx).Find()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	if err = copier.Copy(&resp.Posts, posts); err != nil {
		return nil, err
	}
	if req == nil {
		return
	}

	userId := req.Id
	sysUser, err := dal.SysUserDal.SelectById(l.ctx, userId)
	if err != nil {
		return nil, err
	}
	user := new(types.UserRoles)
	if err = copier.Copy(&user, sysUser); err != nil {
		return nil, err
	}
	resp.User = user

	roleIds, err := dal.SysUserDal.SelectRoleIdsByUserId(l.ctx, userId)
	if err != nil {
		return nil, err
	}
	postIds, err := dal.SysUserDal.SelectPostIdsByUserId(l.ctx, userId)
	if err != nil {
		return nil, err
	}
	resp.RoleIds = roleIds
	resp.PostIds = postIds
	return
}
