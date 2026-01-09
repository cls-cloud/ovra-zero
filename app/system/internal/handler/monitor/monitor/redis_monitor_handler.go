// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package monitor

import (
	"net/http"

	"ovra/app/system/internal/logic/monitor/monitor"
	"ovra/app/system/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
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
