// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package online

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/system/internal/logic/monitor/online"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
)

func OfflineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := online.NewOfflineLogic(r.Context(), svcCtx)
		err := l.Offline(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
