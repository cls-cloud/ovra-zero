package _post

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/system/internal/logic/_post"
	"ovra/app/system/internal/svc"
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
