// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tree

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/demo/internal/logic/demo/tree"
	"ovra/app/demo/internal/svc"
	"ovra/app/demo/internal/types"
)

func AddHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ModifyTreeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := tree.NewAddLogic(r.Context(), svcCtx)
		err := l.Add(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
