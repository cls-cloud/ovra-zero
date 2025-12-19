package auth

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/system/internal/logic/auth"
	"ovra/app/system/internal/svc"
)

func TenantListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewTenantListLogic(r.Context(), svcCtx)
		resp, err := l.TenantList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
