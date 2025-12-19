package dept

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ovra/app/system/internal/logic/dept"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"
)

func UpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bodyMap map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&bodyMap); err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("invalid json: %w", err))
			return
		}
		if val, ok := bodyMap["parentId"]; ok {
			switch v := val.(type) {
			case float64:
				bodyMap["parentId"] = fmt.Sprintf("%.0f", v)
			case int:
				bodyMap["parentId"] = strconv.Itoa(v)
			case string:
			default:
				httpx.ErrorCtx(r.Context(), w, fmt.Errorf("invalid parentId type"))
				return
			}
		}
		data, err := json.Marshal(bodyMap)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("marshal error: %w", err))
			return
		}

		var req types.ModifyDeptReq
		if err := json.Unmarshal(data, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("unmarshal error: %w", err))
			return
		}

		l := dept.NewUpdateLogic(r.Context(), svcCtx)
		err = l.Update(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
