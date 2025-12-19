package plugin

import (
	"ovra/toolkit/auth"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TenantPlugin 多租户插件
type TenantPlugin struct {
	Enabled      bool     // 是否启用多租户
	IgnoreTables []string // 忽略多租户处理的表名
}

// Name 插件名称
func (tp *TenantPlugin) Name() string {
	return "TenantPlugin"
}

// getTenantID 从 context 中获取 tenant_id（string）
func getTenantID(db *gorm.DB) (string, bool) {
	v, ok := db.Statement.Context.Value(auth.TenantIDKey).(string)
	return v, ok
}

// shouldSkip 判断是否跳过当前表的 tenant 限制
func (tp *TenantPlugin) shouldSkip(db *gorm.DB) bool {
	if !tp.Enabled || db.Statement == nil || db.Statement.Table == "" {
		return true
	}
	table := strings.ToLower(db.Statement.Table)
	for _, t := range tp.IgnoreTables {
		if table == strings.ToLower(t) {
			return true
		}
	}
	return false
}

// Initialize 注册 GORM 插件回调
func (tp *TenantPlugin) Initialize(db *gorm.DB) error {
	// Query 钩子
	// Query 钩子：查询时自动添加 tenant_id 条件（若未手动指定）
	if err := db.Callback().Query().Before("gorm:query").
		Register("tenant:query", func(db *gorm.DB) {
			if tp.shouldSkip(db) {
				return
			}
			if tenantID, ok := getTenantID(db); ok {
				// 检查是否已手动指定了 tenant_id 条件
				if whereClause, ok := db.Statement.Clauses["WHERE"]; ok {
					if whereClause.Expression != nil {
						if where, ok := whereClause.Expression.(clause.Where); ok {
							for _, expr := range where.Exprs {
								if eq, ok := expr.(clause.Eq); ok {
									if col, ok := eq.Column.(clause.Column); ok && col.Name == "tenant_id" {
										return // 已设置租户ID，跳过
									}
								}
							}
						}
					}
				}
				db.Statement.AddClause(clause.Where{
					Exprs: []clause.Expression{
						clause.Eq{Column: clause.Column{Table: db.Statement.Table, Name: "tenant_id"}, Value: tenantID},
					},
				})
			}
		}); err != nil {
		return err
	}

	// Create 钩子：创建时自动设置 tenant_id
	if err := db.Callback().Create().Before("gorm:create").
		Register("tenant:create", func(db *gorm.DB) {
			if tp.shouldSkip(db) {
				return
			}
			if tenantID, ok := getTenantID(db); ok {
				if field := db.Statement.Schema.LookUpField("TenantID"); field != nil {
					if db.Statement.ReflectValue.IsValid() && db.Statement.ReflectValue.Kind() == reflect.Struct {
						db.Statement.ReflectValue.FieldByName("TenantID").SetString(tenantID)
					}
				}
			}
		}); err != nil {
		return err
	}

	// Update 钩子
	if err := db.Callback().Update().Before("gorm:update").
		Register("tenant:update", func(db *gorm.DB) {
			if tp.shouldSkip(db) {
				return
			}
			if tenantID, ok := getTenantID(db); ok {
				db.Statement.AddClause(clause.Where{
					Exprs: []clause.Expression{
						clause.Eq{Column: clause.Column{Table: db.Statement.Table, Name: "tenant_id"}, Value: tenantID},
					},
				})
			}
		}); err != nil {
		return err
	}

	// Delete 钩子
	if err := db.Callback().Delete().Before("gorm:delete").
		Register("tenant:delete", func(db *gorm.DB) {
			if tp.shouldSkip(db) {
				return
			}
			if tenantID, ok := getTenantID(db); ok {
				db.Statement.AddClause(clause.Where{
					Exprs: []clause.Expression{
						clause.Eq{Column: clause.Column{Table: db.Statement.Table, Name: "tenant_id"}, Value: tenantID},
					},
				})
			}
		}); err != nil {
		return err
	}

	return nil
}
