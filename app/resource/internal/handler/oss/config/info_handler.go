package config

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/resource/internal/logic/oss/config"
	"ovra/app/resource/internal/svc"
	"ovra/app/resource/internal/types"
)

func InfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := config.NewInfoLogic(r.Context(), svcCtx)
		resp, err := l.Info(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
