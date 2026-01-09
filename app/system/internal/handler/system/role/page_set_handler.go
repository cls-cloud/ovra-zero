// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"net/http"

	"ovra/app/system/internal/logic/system/role"
	"ovra/app/system/internal/svc"
	"ovra/app/system/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PageSetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RolePageSetReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := role.NewPageSetLogic(r.Context(), svcCtx)
		resp, err := l.PageSet(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		rows := make([]map[string]any, 0, len(resp.Rows))
		for _, row := range resp.Rows {
			m := map[string]any{
				"roleId":            row.RoleID,
				"roleName":          row.RoleName,
				"roleKey":           row.RoleKey,
				"roleSort":          row.RoleSort,
				"status":            row.Status,
				"createTime":        row.CreateTime,
				"menuCheckStrictly": row.MenuCheckStrictly,
				"deptCheckStrictly": row.DeptCheckStrictly,
				"dataScope":         row.DataScope,
				"superAdmin":        row.SuperAdmin,
			}

			if row.RoleID == "1" {
				m["roleId"] = 1
			}

			rows = append(rows, m)
		}

		httpx.OkJsonCtx(r.Context(), w, map[string]any{
			"total": resp.Total,
			"rows":  rows,
		})
	}
}
