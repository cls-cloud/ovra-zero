package plugin_test

import (
	"context"
	"os"
	"ovra/toolkit/gorm/plugin"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	slog "log"
)

type User struct {
	ID         int64 `gorm:"primaryKey"`
	Name       string
	TenantID   string
	CreateBy   string
	UpdateBy   string
	CreateTime time.Time
	UpdateTime time.Time
}

func (*User) TableName() string {
	return "test_user"
}

func TestAuditPlugin(t *testing.T) {
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
	err = db.Use(&plugin.AuditPlugin{})
	if err != nil {
		t.Fatalf("failed to use tenant plugin: %v", err)
	}
	// 清空数据库表
	db.Migrator().DropTable(&User{})
	db.AutoMigrate(&User{})

	// 带有 user_id 的 context
	ctx := context.WithValue(context.Background(), "user_id", "tester")

	// 创建记录
	user := User{Name: "Alice"}
	err = db.WithContext(ctx).Create(&user).Error
	assert.NoError(t, err)

	// 查询记录，验证 create_by 和 update_by
	var createdUser User
	err = db.First(&createdUser, user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "tester", createdUser.CreateBy)
	assert.Equal(t, "tester", createdUser.UpdateBy)

	// 更新记录
	err = db.WithContext(ctx).Model(&createdUser).Update("Name", "AliceUpdated").Error
	assert.NoError(t, err)

	// 再次查询
	var updatedUser User
	err = db.First(&updatedUser, user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "AliceUpdated", updatedUser.Name)
	assert.Equal(t, "tester", updatedUser.UpdateBy)
}
