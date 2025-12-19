package plugin_test

import (
	"context"
	slog "log"
	"os"
	"ovra/toolkit/auth"
	"ovra/toolkit/gorm/plugin"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func TestDataScopePlugin(t *testing.T) {
	// 创建数据库连接
	dsn := "root:Pl@1221view@tcp(127.0.0.1:3306)/atlas_zero?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		slog.New(os.Stdout, "\r\n", slog.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:         newLogger,
		NamingStrategy: schema.NamingStrategy{
			// SingularTable: true,
		},
	})
	assert.NoError(t, err)

	// 注册 DataScope 插件
	err = db.Use(&plugin.DataScopePlugin{
		Enabled:      true,
		IgnoreTables: []string{},
	})
	err = db.Use(&plugin.AuditPlugin{})
	assert.NoError(t, err)

	// 清空并建表
	_ = db.Migrator().DropTable(&User{})
	assert.NoError(t, db.AutoMigrate(&User{}))

	t.Run("添加数据", func(t *testing.T) {
		// 插入数据（不同部门）
		ctx := context.Background()
		ctx = context.WithValue(ctx, auth.DataScopeKey, 3)
		ctx = context.WithValue(ctx, auth.UserIDKey, "1")

		err = db.WithContext(ctx).Create(&User{ID: 1, Name: "张三", CreateDept: "100"}).Error
		err = db.WithContext(ctx).Create(&User{ID: 2, Name: "李四", CreateDept: "200"}).Error
		err = db.WithContext(ctx).Create(&User{ID: 3, Name: "王五", CreateDept: "300"}).Error

		if err != nil {
			panic(err)
		}
	})

	t.Run("Scope=3 本部门数据", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, auth.DataScopeKey, 3)
		ctx = context.WithValue(ctx, auth.UserIDKey, "1")
		ctx = context.WithValue(ctx, auth.CurrentDeptKey, "100")
		ctx = context.WithValue(ctx, auth.BellowDeptKey, "100,101")
		ctx = context.WithValue(ctx, auth.CustomerDeptKey, "200,300")

		var users []User
		err := db.WithContext(ctx).Find(&users).Error
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, int64(1), users[0].ID)
	})

	t.Run("Scope=5 仅本人数据", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, auth.DataScopeKey, 5)
		ctx = context.WithValue(ctx, auth.UserIDKey, "1")

		var users []User
		err := db.WithContext(ctx).Where("id = ?", 1).Find(&users).Error
		assert.NoError(t, err)
		assert.Len(t, users, 1)
	})

	t.Run("Scope=4 本部门及以下", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, auth.DataScopeKey, 4)
		ctx = context.WithValue(ctx, auth.BellowDeptKey, "100,101,200")

		var users []User
		err := db.WithContext(ctx).Find(&users).Error
		assert.NoError(t, err)
		// 应该查出 1（100）、2（200）
		assert.Len(t, users, 2)
	})
}
