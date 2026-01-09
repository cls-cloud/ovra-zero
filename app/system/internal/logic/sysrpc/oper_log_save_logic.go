package sysrpclogic

import (
	"context"
	"ovra/app/system/internal/dal/model"
	"ovra/toolkit/utils"
	"time"

	"ovra/app/system/internal/svc"
	"ovra/app/system/pb/system"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOperLogSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogSaveLogic {
	return &OperLogSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OperLogSaveLogic) OperLogSave(in *system.OperLogReq) (*system.EmptyResp, error) {
	q := l.svcCtx.Dal.Query
	id := utils.GetID()
	err := q.SysOperLog.WithContext(l.ctx).Create(&model.SysOperLog{
		OperID:        id,
		TenantID:      in.TenantId,
		Title:         in.Title,
		BusinessType:  in.BusinessType,
		Method:        in.Method,
		RequestMethod: in.RequestMethod,
		OperatorType:  in.OperatorType,
		OperName:      in.OperName,
		DeptName:      in.DeptName,
		OperURL:       in.OperUrl,
		OperIP:        in.OperIp,
		OperLocation:  in.OperLocation,
		OperParam:     in.OperParam,
		JSONResult:    in.JsonResult,
		Status:        in.Status,
		ErrorMsg:      in.ErrorMsg,
		OperTime:      time.Now(),
		CostTime:      in.CostTime,
	})
	if err != nil {
		return nil, err
	}
	return &system.EmptyResp{}, nil
}
