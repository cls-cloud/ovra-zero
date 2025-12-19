package menu

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/system/internal/logic/menu"
	"ovra/app/system/internal/svc"
)

func GetRoutersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewGetRoutersLogic(r.Context(), svcCtx)
		resp, err := l.GetRouters()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
