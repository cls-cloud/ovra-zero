package oss

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/resource/internal/logic/oss"
	"ovra/app/resource/internal/svc"
)

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := oss.NewUploadLogic(r.Context(), svcCtx, r)
		err := l.Upload()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
