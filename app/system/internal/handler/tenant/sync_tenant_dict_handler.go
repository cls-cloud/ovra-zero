package tenant

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/system/internal/logic/tenant"
	"ovra/app/system/internal/svc"
)

func SyncTenantDictHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := tenant.NewSyncTenantDictLogic(r.Context(), svcCtx)
		err := l.SyncTenantDict()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
