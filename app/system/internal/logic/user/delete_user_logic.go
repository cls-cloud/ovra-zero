package user

import (
	"context"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.CodeReq) error {
	userIds := strings.Split(req.Code, ",")
	if len(userIds) != 0 {
		if err := l.svcCtx.Dal.SysUserDal.DeleteBatch(l.ctx, userIds); err != nil {
			return err
		}
	}
	return nil
}
