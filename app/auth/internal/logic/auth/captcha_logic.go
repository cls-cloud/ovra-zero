// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"ovra/app/auth/internal/svc"
	"ovra/app/auth/internal/types"

	"github.com/mojocn/base64Captcha"
	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaLogic {
	return &CaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var Store = base64Captcha.DefaultMemStore

func (l *CaptchaLogic) Captcha() (resp *types.CaptchaResp, err error) {
	driver := base64Captcha.NewDriverDigit(80, 250, 4, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, Store)
	id, b64s, answer, err := cp.Generate()
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.Rds.SetexCtx(l.ctx, "captcha:"+id, answer, 5*60)
	if err != nil {
		return nil, err
	}
	return &types.CaptchaResp{
		Img:            b64s[22:],
		Uuid:           id,
		CaptchaEnabled: l.svcCtx.Config.Captcha.Enabled,
	}, nil
}
