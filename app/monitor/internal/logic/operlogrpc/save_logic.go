package operlogrpclogic

import (
	"context"
	"ovra/app/monitor/internal/dal/model"
	"ovra/toolkit/utils"
	"time"

	"ovra/app/monitor/internal/svc"
	"ovra/app/monitor/pb/monitor"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveLogic {
	return &SaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SaveLogic) Save(in *monitor.OperLogReq) (*monitor.EmptyResp, error) {
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
	return &monitor.EmptyResp{}, err
}
