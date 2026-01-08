package plugin

import (
	"ovra/toolkit/auth"
	"reflect"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type AuditPlugin struct{}

func (ap *AuditPlugin) Name() string {
	return "AuditPlugin"
}

func getUserID(db *gorm.DB) (string, bool) {
	v, ok := db.Statement.Context.Value(auth.UserIDKey).(string)
	return v, ok
}

func (ap *AuditPlugin) Initialize(db *gorm.DB) error {
	// ===== Create =====
	if err := db.Callback().Create().Before("gorm:create").
		Register("audit:create", func(db *gorm.DB) {
			userID, ok := getUserID(db)
			if !ok || db.Statement == nil || db.Statement.Schema == nil {
				return
			}
			now := time.Now()
			WalkStruct(db.Statement.ReflectValue, func(v reflect.Value) {
				processAuditCreate(db, v, userID, now)
			})
		}); err != nil {
		return err
	}
	// ===== Update =====
	if err := db.Callback().Update().Before("gorm:update").
		Register("audit:update", func(db *gorm.DB) {
			userID, ok := getUserID(db)
			if !ok {
				return
			}
			now := time.Now()
			db.Statement.SetColumn("UpdateBy", userID)
			db.Statement.SetColumn("UpdateTime", now)
		}); err != nil {
		return err
	}
	return nil
}

func processAuditCreate(db *gorm.DB, v reflect.Value, userID string, now time.Time) {
	// create_by
	if field := db.Statement.Schema.LookUpField("CreateBy"); field != nil {
		if f := v.FieldByName("CreateBy"); f.IsValid() && f.CanSet() && f.IsZero() {
			setUserValue(f, userID)
		}
	}
	// update_by
	if field := db.Statement.Schema.LookUpField("UpdateBy"); field != nil {
		if f := v.FieldByName("UpdateBy"); f.IsValid() && f.CanSet() && f.IsZero() {
			setUserValue(f, userID)
		}
	}
	// create_time
	if field := db.Statement.Schema.LookUpField("CreateTime"); field != nil {
		if f := v.FieldByName("CreateTime"); f.IsValid() && f.CanSet() {
			if f.Type() == reflect.TypeOf(time.Time{}) {
				f.Set(reflect.ValueOf(now))
			}
		}
	}
	// update_time
	if field := db.Statement.Schema.LookUpField("UpdateTime"); field != nil {
		if f := v.FieldByName("UpdateTime"); f.IsValid() && f.CanSet() {
			if f.Type() == reflect.TypeOf(time.Time{}) {
				f.Set(reflect.ValueOf(now))
			}
		}
	}
}

func setUserValue(f reflect.Value, userID string) {
	switch f.Kind() {
	case reflect.String:
		f.SetString(userID)
	case reflect.Int, reflect.Int64:
		if id, err := strconv.ParseInt(userID, 10, 64); err == nil {
			f.SetInt(id)
		}
	}
}
