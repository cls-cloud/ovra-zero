package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/system/internal/logic/user"
	"ovra/app/system/internal/svc"
)

func QueryUserDetailInsertHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewQueryUserDetailInsertLogic(r.Context(), svcCtx)
		resp, err := l.QueryUserDetailInsert()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
