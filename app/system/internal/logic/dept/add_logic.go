package dept

import (
	"context"
	"ovra/app/system/internal/dal/model"
	"ovra/toolkit/errx"
	"ovra/toolkit/utils"

	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req *types.ModifyDeptReq) error {
	q := l.svcCtx.Dal.Query
	var ancestors string
	parentId := req.ParentID
	//获取祖级列表
	if parentId != "0" {
		if err := q.SysDept.WithContext(l.ctx).Select(q.SysDept.Ancestors).Where(q.SysDept.DeptID.Eq(parentId)).Scan(&ancestors); err != nil {
			return errx.GORMErr(err)
		}
		ancestors = ancestors + "," + req.ParentID
	} else {
		ancestors = req.ParentID
	}
	deptId := utils.GetID()
	if req.DeptID != "" {
		deptId = req.DeptID
	}
	dept := &model.SysDept{
		DeptID:    deptId,
		DeptName:  req.DeptName,
		Ancestors: ancestors,
		ParentID:  parentId,
		OrderNum:  req.OrderNum,
		Phone:     req.Phone,
		Email:     req.Email,
		Status:    req.Status,
	}
	if req.TenantID != "" {
		dept.TenantID = req.TenantID
	}
	if err := q.SysDept.WithContext(l.ctx).Create(dept); err != nil {
		return errx.GORMErr(err)
	}
	return nil
}
