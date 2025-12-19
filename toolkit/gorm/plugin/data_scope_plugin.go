package plugin

//1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限 5：仅本人数据权限 6：部门及以下或本人数据权限

import (
	"fmt"
	"ovra/toolkit/auth"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DataScopePlugin 数据权限插件
type DataScopePlugin struct {
	Enabled      bool
	IgnoreTables []string
}

// 插件名称
func (dsp *DataScopePlugin) Name() string {
	return "DataScopePlugin"
}

// 注册插件
func (dsp *DataScopePlugin) Initialize(db *gorm.DB) error {
	return db.Callback().Query().Before("gorm:query").
		Register("data_scope:query", func(db *gorm.DB) {
			if dsp.shouldSkip(db) {
				return
			}

			dataScope, _ := db.Statement.Context.Value(auth.DataScopeKey).(int)
			userID, _ := db.Statement.Context.Value(auth.UserIDKey).(string)
			deptID, _ := db.Statement.Context.Value(auth.CurrentDeptKey).(string)
			bellowDeptID, _ := db.Statement.Context.Value(auth.BellowDeptKey).(string)
			customerDeptID, _ := db.Statement.Context.Value(auth.CustomerDeptKey).(string)

			sql := dsp.genWhereSQL(db.Statement.Table, dataScope, userID, deptID, bellowDeptID, customerDeptID)
			if sql != "" {
				db.Statement.AddClause(clause.Where{
					Exprs: []clause.Expression{
						clause.Expr{SQL: sql},
					},
				})
			}
		})
}

// 判断是否跳过
func (dsp *DataScopePlugin) shouldSkip(db *gorm.DB) bool {
	if !dsp.Enabled || db.Statement == nil || db.Statement.Table == "" {
		return true
	}
	table := strings.ToLower(db.Statement.Table)
	for _, t := range dsp.IgnoreTables {
		if table == strings.ToLower(t) {
			return true
		}
	}
	return false
}

func formatINList(csv string) string {
	ids := strings.Split(csv, ",")
	for i := range ids {
		ids[i] = fmt.Sprintf("'%s'", strings.TrimSpace(ids[i]))
	}
	return fmt.Sprintf("(%s)", strings.Join(ids, ","))
}

func (dsp *DataScopePlugin) genWhereSQL(table string, scope int, userID, deptID, bellowDeptID, customerDeptID string) string {
	colPrefix := table
	if colPrefix != "" {
		colPrefix += "."
	}

	switch scope {
	case 1:
		return ""
	case 2:
		return fmt.Sprintf("%screate_dept IN %s", colPrefix, formatINList(customerDeptID))
	case 3:
		return fmt.Sprintf("%screate_dept = '%s'", colPrefix, deptID)
	case 4:
		return fmt.Sprintf("%screate_dept IN %s", colPrefix, formatINList(bellowDeptID))
	case 5:
		return fmt.Sprintf("%screate_by = '%s'", colPrefix, userID)
	case 6:
		return fmt.Sprintf("(%screate_dept IN %s OR %screate_by = '%s')", colPrefix, formatINList(bellowDeptID), colPrefix, userID)
	default:
		return fmt.Sprintf("%screate_by = '%s'", colPrefix, userID)
	}
}
