// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package _post

import (
	"net/http"

	"ovra/app/system/internal/logic/system/_post"
	"ovra/app/system/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetDeptTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := _post.NewGetDeptTreeLogic(r.Context(), svcCtx)
		resp, err := l.GetDeptTree()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
