// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"ovra/app/system/internal/logic/system/menu"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AddHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 兼容前端处理
		body, err := io.ReadAll(r.Body)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		var raw map[string]any
		if err := json.Unmarshal(body, &raw); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if v, ok := raw["parentId"]; ok {
			switch val := v.(type) {
			case string:
				raw["parentId"] = val
			case float64:
				raw["parentId"] = strconv.FormatFloat(val, 'f', -1, 64)
			case int:
				raw["parentId"] = strconv.Itoa(val)
			}
		}
		var req types.ModifyMenuReq
		bs, _ := json.Marshal(raw)
		if err := json.Unmarshal(bs, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := menu.NewAddLogic(r.Context(), svcCtx)
		if err := l.Add(&req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
