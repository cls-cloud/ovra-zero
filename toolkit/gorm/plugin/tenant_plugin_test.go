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
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func TestTenantPlugin(t *testing.T) {
	// 创建数据库连接
	dsn := "root:Pl@1221view@tcp(127.0.0.1:3306)/atlas_zero?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		slog.New(os.Stdout, "\r\n", slog.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢查询sql阀值
			Colorful:      true,        //禁用彩色打印
			LogLevel:      logger.Info,
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{

		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy:                           schema.NamingStrategy{
			//SingularTable: true, //表名是否加s
		},
	})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// 使用 TenantPlugin 插件
	err = db.Use(&plugin.TenantPlugin{
		Enabled:      true,
		IgnoreTables: []string{},
	})
	if err != nil {
		t.Fatalf("failed to use tenant plugin: %v", err)
	}

	// 清空数据库表（确保每次测试是干净的）
	db.Migrator().DropTable(&User{})
	db.AutoMigrate(&User{})
	//
	//// 测试：没有 tenant_id 的情况下插入
	//t.Run("No TenantID", func(t *testing.T) {
	//	// 没有传递 tenant_id，插入时应该不会插入 tenant_id 字段
	//	err := db.Create(&User{Name: "John Doe"}).Error
	//	assert.NoError(t, err)
	//
	//	// 查询数据，应该没有 tenant_id
	//	var user User
	//	db.First(&user)
	//	assert.Equal(t, int64(0), user.TenantID) // 没有 tenant_id
	//})
	//
	//// 测试：传递 tenant_id 进行插入
	//t.Run("With TenantID", func(t *testing.T) {
	//	// 使用 context 传递 tenant_id
	//	tenantID := int64(1)
	//	db = db.WithContext(context.WithValue(context.Background(), "tenant_id", tenantID))
	//
	//	// 插入数据时，应该自动设置 tenant_id
	//	err := db.Create(&User{Name: "Jane Doe"}).Error
	//	assert.NoError(t, err)
	//
	//	// 查询数据，应该有 tenant_id 设置为 1
	//	var user User
	//	db.First(&user)
	//	assert.Equal(t, tenantID, user.TenantID) // tenant_id 应该是 1
	//})

	// 测试：查询时自动加上 tenant_id 过滤条件
	t.Run("Query with TenantID", func(t *testing.T) {
		// 使用 context 传递 tenant_id
		tenantID := "1"
		db = db.WithContext(context.WithValue(context.Background(), auth.TenantIDKey, tenantID))

		// 插入数据时，应该自动设置 tenant_id
		err := db.Create(&User{Name: "Alice Doe"}).Error
		assert.NoError(t, err)

		// 查询数据
		var user User
		err = db.Find(&user).Error
		logx.Info("user", user)
		assert.NoError(t, err)
		assert.Equal(t, tenantID, user.TenantID) // tenant_id 应该是 1
		//assert.Equal(t, int64(0), user.TenantID)
	})
}
