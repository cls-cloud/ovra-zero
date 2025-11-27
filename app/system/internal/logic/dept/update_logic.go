package dept

import (
	"context"
	"toolkit/errx"
	"toolkit/utils"

	"system/internal/svc"
	"system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.ModifyDeptReq) error {
	toMapOmit := utils.StructToMapOmit(req.DeptBase, nil, []string{"Children", "CreateTime"}, true)
	if req.DeptCategory == "" {
		toMapOmit["dept_category"] = req.DeptCategory
	}
	if req.Phone == "" {
		toMapOmit["phone"] = req.Phone
	}
	if req.Email == "" {
		toMapOmit["email"] = req.Email
	}
	if _, err := l.svcCtx.Dal.Query.SysDept.WithContext(l.ctx).Where(l.svcCtx.Dal.Query.SysDept.DeptID.Eq(req.DeptID)).Updates(toMapOmit); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
