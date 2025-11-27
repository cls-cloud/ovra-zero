package dept

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"system/internal/svc"
	"system/internal/types"
	"toolkit/errx"

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
	q := l.svcCtx.Dal.Query
	sysDept, err := q.SysDept.WithContext(l.ctx).Where(q.SysDept.DeptID.Neq(req.Id)).
		Where(q.SysDept.Ancestors.NotRegexp(fmt.Sprintf("^%d", req.Id))).
		Find()
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	if err = copier.Copy(&resp, sysDept); err != nil {
		return nil, err
	}
	return
}
