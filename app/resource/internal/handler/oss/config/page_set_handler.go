package config

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/resource/internal/logic/oss/config"
	"ovra/app/resource/internal/svc"
	"ovra/app/resource/internal/types"
)

func PageSetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageSetOssConfigReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := config.NewPageSetLogic(r.Context(), svcCtx)
		resp, err := l.PageSet(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
