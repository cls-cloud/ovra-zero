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

func PageSetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageSetDemoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := demo.NewPageSetLogic(r.Context(), svcCtx)
		resp, err := l.PageSet(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
