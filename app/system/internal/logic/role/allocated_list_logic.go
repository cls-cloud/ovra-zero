package role

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"toolkit/errx"

	"system/internal/svc"
	"system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllocatedListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAllocatedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllocatedListLogic {
	return &AllocatedListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllocatedListLogic) AllocatedList(req *types.AllocatedReq) (resp *types.AllocatedResp, err error) {
	data, err := l.GetUserList(req, true)
	resp = new(types.AllocatedResp)
	resp = data
	return
}

func (l *AllocatedListLogic) GetUserList(req *types.AllocatedReq, type_ bool) (resp *types.AllocatedResp, err error) {
	resp = new(types.AllocatedResp)
	offset := (req.PageNum - 1) * req.PageSize
	d := l.svcCtx.Dal.Query
	var userIds []string
	if err = d.SysUserRole.WithContext(l.ctx).Select(d.SysUserRole.UserID).Where(d.SysUserRole.RoleID.Eq(req.RoleId)).Scan(&userIds); err != nil {
		return nil, errx.GORMErr(err)
	}
	do := d.SysUser.WithContext(l.ctx)
	if req.UserName != "" {
		do = do.Where(d.SysUser.UserName.Like(fmt.Sprintf("%%%s%%", req.UserName)))
	}
	if req.PhoneNumber != "" {
		do = do.Where(d.SysUser.Phonenumber.Like(fmt.Sprintf("%%%s%%", req.PhoneNumber)))
	}
	if type_ {
		do = do.Where(d.SysUser.UserID.In(userIds...))
	} else {
		do = do.Where(d.SysUser.UserID.NotIn(userIds...))
	}
	users, count, err := do.FindByPage(int(offset), int(req.PageSize))
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp.Rows = make([]types.UserBase, len(users))
	for i, row := range users {
		var user types.UserBase
		_ = copier.Copy(&user, row)
		resp.Rows[i] = user
	}
	resp.Total = count
	return
}
