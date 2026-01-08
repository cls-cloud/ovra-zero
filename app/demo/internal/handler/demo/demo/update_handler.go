// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package demo

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/demo/internal/logic/demo/demo"
	"ovra/app/demo/internal/svc"
	"ovra/app/demo/internal/types"
)

func UpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ModifyDemoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := demo.NewUpdateLogic(r.Context(), svcCtx)
		err := l.Update(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
