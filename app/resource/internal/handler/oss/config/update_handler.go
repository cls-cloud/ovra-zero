package config

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/resource/internal/logic/oss/config"
	"ovra/app/resource/internal/svc"
	"ovra/app/resource/internal/types"
)

func UpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ModifyOssConfigReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := config.NewUpdateLogic(r.Context(), svcCtx)
		err := l.Update(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
