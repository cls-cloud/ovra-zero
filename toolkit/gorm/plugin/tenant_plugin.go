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
	// ===== Query =====
	if err := db.Callback().Query().Before("gorm:query").
		Register("tenant:query", func(db *gorm.DB) {
			if tp.shouldSkip(db) {
				return
			}
			if tenantID, ok := getTenantID(db); ok {
				addTenantWhereIfAbsent(db, tenantID)
			}
		}); err != nil {
		return err
	}

	// ===== Create =====
	if err := db.Callback().Create().Before("gorm:create").
		Register("tenant:create", func(db *gorm.DB) {
			if tp.shouldSkip(db) {
				return
			}
			tenantID, ok := getTenantID(db)
			if !ok || tenantID == "" {
				return
			}
			WalkStruct(db.Statement.ReflectValue, func(v reflect.Value) {
				setTenantIDIfEmpty(v, tenantID)
			})
		}); err != nil {
		return err
	}
	// ===== Update =====
	if err := db.Callback().Update().Before("gorm:update").
		Register("tenant:update", func(db *gorm.DB) {
			if tp.shouldSkip(db) {
				return
			}
			if tenantID, ok := getTenantID(db); ok {
				addTenantWhereIfAbsent(db, tenantID)
			}
		}); err != nil {
		return err
	}
	// ===== Delete =====
	if err := db.Callback().Delete().Before("gorm:delete").
		Register("tenant:delete", func(db *gorm.DB) {
			if tp.shouldSkip(db) {
				return
			}
			if tenantID, ok := getTenantID(db); ok {
				addTenantWhereIfAbsent(db, tenantID)
			}
		}); err != nil {
		return err
	}
	return nil
}

func setTenantIDIfEmpty(v reflect.Value, tenantID string) {
	field := v.FieldByName("TenantID")
	if !field.IsValid() || !field.CanSet() {
		return
	}
	if field.Kind() != reflect.String {
		return
	}
	if field.String() == "" {
		field.SetString(tenantID)
	}
}

func addTenantWhereIfAbsent(db *gorm.DB, tenantID string) {
	if _, ok := db.Statement.Clauses["tenant:where"]; ok {
		return
	}
	db.Statement.AddClause(clause.Where{
		Exprs: []clause.Expression{
			clause.Eq{
				Column: clause.Column{
					Table: db.Statement.Table,
					Name:  "tenant_id",
				},
				Value: tenantID,
			},
		},
	})
	db.Statement.Settings.Store("tenant:where", true)
}
