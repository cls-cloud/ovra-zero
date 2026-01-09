package role

import (
	"context"
	"fmt"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageSetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPageSetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageSetLogic {
	return &PageSetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageSetLogic) PageSet(req *types.RolePageSetReq) (resp *types.RolePageSetResp, err error) {
	resp = new(types.RolePageSetResp)
	offset := (req.PageNum - 1) * req.PageSize
	q := l.svcCtx.Dal.Query
	do := q.SysRole.WithContext(l.ctx)
	if req.RoleName != "" {
		do = do.Where(q.SysRole.RoleName.Like(fmt.Sprintf("%%%s%%", req.RoleName)))
	}
	if req.RoleKey != "" {
		do = do.Where(q.SysRole.RoleKey.Like(fmt.Sprintf("%%%s%%", req.RoleKey)))
	}
	if req.Status != "" {
		do = do.Where(q.SysRole.Status.Eq(req.Status))
	}
	if req.BeginTime != "" {
		beginTime, _ := time.Parse(time.DateTime, req.BeginTime)
		do = do.Where(q.SysRole.CreateTime.Gte(beginTime))
	}
	if req.EndTime != "" {
		endTime, _ := time.Parse(time.DateTime, req.EndTime)
		do = do.Where(q.SysRole.CreateTime.Lte(endTime))
	}
	result, total, err := do.Order(q.SysRole.RoleSort, q.SysRole.CreateTime.Desc()).FindByPage(int(offset), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	roles := make([]types.RoleBase, 0)
	for _, v := range result {
		roles = append(roles, types.RoleBase{
			RoleID:            v.RoleID,
			RoleName:          v.RoleName,
			RoleKey:           v.RoleKey,
			RoleSort:          v.RoleSort,
			Status:            v.Status,
			CreateTime:        v.CreateTime.Format(time.DateTime),
			MenuCheckStrictly: v.MenuCheckStrictly,
			DeptCheckStrictly: v.DeptCheckStrictly,
			DataScope:         v.DataScope,
			SuperAdmin:        v.RoleID == "1",
		})
	}
	resp.Rows = roles
	resp.Total = total
	return
}
