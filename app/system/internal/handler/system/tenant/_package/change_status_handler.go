// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package _package

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/system/internal/logic/system/tenant/_package"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
)

func ChangeStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangeStatusTenantPackageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := _package.NewChangeStatusLogic(r.Context(), svcCtx)
		err := l.ChangeStatus(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
