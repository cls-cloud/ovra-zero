package role

import (
	"context"
	"system/internal/logic/user"
	"system/internal/svc"
	"system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeptTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptTreeLogic {
	return &DeptTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeptTreeLogic) DeptTree(req *types.IdReq) (resp *types.DeptTreeResp, err error) {
	resp = new(types.DeptTreeResp)
	tree, err := user.NewGetDeptTreeLogic(l.ctx, l.svcCtx).GetDeptTree()
	resp.Depts = tree
	//获取角色已经分配的部门
	q := l.svcCtx.Dal.Query
	roleDept, err := q.SysRoleDept.WithContext(l.ctx).Where(q.SysRoleDept.RoleID.Eq(req.Id)).Find()
	deptIds := make([]string, 0, len(roleDept))
	for i := range roleDept {
		deptIds[i] = roleDept[i].DeptID
	}
	resp.CheckedKeys = deptIds
	return
}
