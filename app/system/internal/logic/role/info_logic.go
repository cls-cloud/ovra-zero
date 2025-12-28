package role

import (
	"context"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
	"ovra/toolkit/errx"

	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoLogic) Info(req *types.IdReq) (resp *types.RoleBase, err error) {
	resp = new(types.RoleBase)
	sysRole, err := l.svcCtx.Dal.SysRoleDal.SelectById(l.ctx, req.Id)
	if err != nil {
		return nil, errx.GORMErr(err)
	}
	if err = copier.Copy(&resp, sysRole); err != nil {
		return nil, err
	}
	return
}
