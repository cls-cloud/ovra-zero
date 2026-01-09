// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logininfor

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/system/internal/logic/monitor/logininfor"
	"ovra/app/system/internal/svc"
)

func CleanHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logininfor.NewCleanLogic(r.Context(), svcCtx)
		err := l.Clean()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
