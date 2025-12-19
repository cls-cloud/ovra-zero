package auth

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/system/internal/logic/auth"
	"ovra/app/system/internal/svc"
)

func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewLogoutLogic(r.Context(), svcCtx, r)
		err := l.Logout()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
