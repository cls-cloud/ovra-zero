package operLog

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/monitor/internal/logic/operLog"
	"ovra/app/monitor/internal/svc"
)

func CleanHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := operLog.NewCleanLogic(r.Context(), svcCtx)
		err := l.Clean()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
