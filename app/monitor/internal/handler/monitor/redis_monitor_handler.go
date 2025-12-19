package monitor

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/monitor/internal/logic/monitor"
	"ovra/app/monitor/internal/svc"
)

func RedisMonitorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := monitor.NewRedisMonitorLogic(r.Context(), svcCtx)
		resp, err := l.RedisMonitor()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
