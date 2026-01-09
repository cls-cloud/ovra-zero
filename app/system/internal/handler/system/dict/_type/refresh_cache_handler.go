// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package _type

import (
	"net/http"

	"ovra/app/system/internal/logic/system/dict/_type"
	"ovra/app/system/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RefreshCacheHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := _type.NewRefreshCacheLogic(r.Context(), svcCtx)
		err := l.RefreshCache()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
