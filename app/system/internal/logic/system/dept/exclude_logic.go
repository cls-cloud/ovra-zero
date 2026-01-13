package dept

import (
	"context"
	"fmt"
	"ovra/app/system/internal/dal/model"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"ovra/toolkit/errx"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExcludeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExcludeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExcludeLogic {
	return &ExcludeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExcludeLogic) Exclude(req *types.IdReq) (resp []*types.DeptBase, err error) {
	db := l.svcCtx.Dal.Db
	switch db.Dialector.Name() {
	case "mysql":
		return l.excludeMysql(req)
	case "postgres":
		return l.excludePgsql(req)
	default:
		return nil, fmt.Errorf("unsupported database: %s", db.Dialector.Name())
	}
}

func (l *ExcludeLogic) excludeMysql(req *types.IdReq) (resp []*types.DeptBase, err error) {
	q := l.svcCtx.Dal.Query
	list, err := q.SysDept.WithContext(l.ctx).Where(q.SysDept.DeptID.Neq(req.Id)).
		Where(q.SysDept.Ancestors.NotRegexp(fmt.Sprintf("^%s", req.Id))).
		Find()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = l.toDeptBaseList(list)
	return
}

func (l *ExcludeLogic) excludePgsql(req *types.IdReq) (resp []*types.DeptBase, err error) {
	db := l.svcCtx.Query.SysDept.
		WithContext(l.ctx).
		UnderlyingDB()
	var list []*model.SysDept
	sql := `
		SELECT *
		FROM sys_dept
		WHERE dept_id <> ?
		  AND (
		        ancestors IS NULL
		        OR ancestors !~ ?
		      )
	`
	pattern := fmt.Sprintf("(^|,)%s(,|$)", req.Id)

	if err = db.Raw(sql, req.Id, pattern).Scan(&list).Error; err != nil {
		return nil, errx.GORMErr(err)
	}
	resp = l.toDeptBaseList(list)
	return
}

func (l *ExcludeLogic) toDeptBaseList(list []*model.SysDept) []*types.DeptBase {
	resp := make([]*types.DeptBase, 0, len(list))

	for _, r := range list {
		leader := r.Leader // 防止以后结构变化
		resp = append(resp, &types.DeptBase{
			DeptID:       r.DeptID,
			ParentID:     r.ParentID,
			Ancestors:    r.Ancestors,
			DeptName:     r.DeptName,
			DeptCategory: r.DeptCategory,
			OrderNum:     r.OrderNum,
			Leader:       &leader,
			Phone:        r.Phone,
			Email:        r.Email,
			Status:       r.Status,
			CreateTime:   r.CreateTime.Format("2006-01-02 15:04:05"),
			Children:     nil,
		})
	}

	return resp
}
