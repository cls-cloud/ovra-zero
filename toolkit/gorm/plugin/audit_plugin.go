package plugin

import (
	"ovra/toolkit/auth"
	"reflect"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// AuditPlugin（用于自动填充 create_by 和 update_by）
type AuditPlugin struct{}

// Name 插件名称
func (ap *AuditPlugin) Name() string {
	return "AuditPlugin"
}

// getUserID 从 context 中获取 user_id（string 类型）
func getUserID(db *gorm.DB) (string, bool) {
	v, ok := db.Statement.Context.Value(auth.UserIDKey).(string)
	return v, ok
}

// Initialize 注册 GORM 插件回调
func (ap *AuditPlugin) Initialize(db *gorm.DB) error {
	// 创建时设置 create_by 和 update_by
	if err := db.Callback().Create().Before("gorm:create").
		Register("audit:create", func(db *gorm.DB) {
			if userID, ok := getUserID(db); ok && db.Statement != nil && db.Statement.Schema != nil {
				val := db.Statement.ReflectValue

				// 如果是指针，先 Elem() 一下
				if val.Kind() == reflect.Ptr {
					val = val.Elem()
				}

				switch val.Kind() {
				case reflect.Struct:
					processStruct(db, val, userID)
				case reflect.Slice:
					for i := 0; i < val.Len(); i++ {
						elem := val.Index(i)
						if elem.Kind() == reflect.Ptr {
							elem = elem.Elem()
						}
						if elem.Kind() == reflect.Struct {
							processStruct(db, elem, userID)
						}
					}
				}

			}
		}); err != nil {
		return err
	}

	// 更新时设置 update_by
	if err := db.Callback().Update().Before("gorm:update").
		Register("audit:update", func(db *gorm.DB) {
			if userID, ok := getUserID(db); ok && db.Statement != nil && db.Statement.Schema != nil {
				now := time.Now()
				// 设置 update_by 字段（兼容结构体更新或 map 更新）
				db.Statement.SetColumn("UpdateBy", userID)

				if field := db.Statement.Schema.LookUpField("UpdateTime"); field != nil {
					if f := db.Statement.ReflectValue.FieldByName("UpdateTime"); f.IsValid() && f.CanSet() {
						if f.Type() == reflect.TypeOf(time.Time{}) {
							f.Set(reflect.ValueOf(now))
						}
					}
				}
			}

		}); err != nil {
		return err
	}

	return nil
}

func processStruct(db *gorm.DB, val reflect.Value, userID string) {
	now := time.Now()
	if val.Kind() == reflect.Struct {
		// 设置 create_by
		if field := db.Statement.Schema.LookUpField("CreateBy"); field != nil {
			if f := val.FieldByName("CreateBy"); f.IsValid() && f.CanSet() {
				switch f.Type().Kind() {
				case reflect.String:
					f.SetString(userID)
				case reflect.Int, reflect.Int64:
					if id, err := strconv.ParseInt(userID, 10, 64); err == nil {
						f.SetInt(id)
					}
				}
			}
		}
		// 设置 update_by
		if field := db.Statement.Schema.LookUpField("UpdateBy"); field != nil {
			if f := val.FieldByName("UpdateBy"); f.IsValid() && f.CanSet() {
				switch f.Type().Kind() {
				case reflect.String:
					f.SetString(userID)
				case reflect.Int, reflect.Int64:
					if id, err := strconv.ParseInt(userID, 10, 64); err == nil {
						f.SetInt(id)
					}
				}
			}
		}
		if field := db.Statement.Schema.LookUpField("CreateTime"); field != nil {
			if f := val.FieldByName("CreateTime"); f.IsValid() && f.CanSet() {
				if f.Type() == reflect.TypeOf(time.Time{}) {
					f.Set(reflect.ValueOf(now))
				}
			}
		}
		if field := db.Statement.Schema.LookUpField("UpdateTime"); field != nil {
			if f := val.FieldByName("UpdateTime"); f.IsValid() && f.CanSet() {
				if f.Type() == reflect.TypeOf(time.Time{}) {
					f.Set(reflect.ValueOf(now))
				}
			}
		}
	}
}
