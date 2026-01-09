// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package oss

import (
	"net/http"

	"ovra/app/system/internal/logic/resource/oss"
	"ovra/app/system/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
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
