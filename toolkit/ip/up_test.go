package ip

import (
	"fmt"
	"ovra/toolkit/gorm/plugin"
	"testing"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Logininfor struct {
	InfoID        int    `gorm:"column:info_id;primaryKey;autoIncrement:true" json:"info_id"`
	IPAddr        string `gorm:"column:ipaddr" json:"ipaddr"`
	LoginLocation string `gorm:"column:login_location" json:"login_location"`
}

func TestA(t *testing.T) {
	db := NewDb()
	var logininfor []Logininfor
	if err := db.Raw("select info_id, ipaddr, login_location from sys_logininfor").Scan(&logininfor).Error; err != nil {
		fmt.Println("Error:", err)
	}
	for _, l := range logininfor {
		// 查询公网 IP 归属地信息
		info, err := LookupIP(l.IPAddr)
		if err != nil {
			logx.Info("查询失败")
			continue
		}
		fmt.Printf("IP: %s\n国家: %s\n省份: %s\n城市: %s\n运营商: %s\n经纬度: %s,%s\n",
			info.IP, info.Country, info.Region, info.City, info.ISP, info.Lat, info.Lng)
		location := fmt.Sprintf("%s|%s|%s|%s", info.Country, info.Region, info.City, info.ISP)
		if err := db.Exec("update sys_logininfor set login_location = ? where info_id = ?", location, l.InfoID).Error; err != nil {
			fmt.Println("更新失败:", err)
		}
		fmt.Println("更新成功")
	}

}

func NewDb() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		"root", "Pl@1221view", "127.0.0.1", 3306, "atlas_zero",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy:                           schema.NamingStrategy{},
	})
	if err != nil {
		panic(err)
	}

	if err = db.Use(&plugin.AuditPlugin{}); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to MySQL database.")
	return db
}
