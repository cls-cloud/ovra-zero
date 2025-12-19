package ip

import (
	"fmt"
	"testing"

	"github.com/zeromicro/go-zero/core/logx"
)

func TestIsPrivateIP(t *testing.T) {
	// 示例 IP
	testIPs := []string{
		"192.168.1.1",     // 内网
		"10.0.0.5",        // 内网
		"8.8.8.8",         // 公网
		"240e:0c::1",      // 公网 IPv6
		"::1",             // 回环地址
		"114.114.114.114", // 公网
	}

	for _, ipStr := range testIPs {
		fmt.Println("==============")
		fmt.Println("查询 IP：", ipStr)

		if IsPrivateIP(ipStr) {
			fmt.Println("这是内网地址或回环地址，无需查询归属地。")
			continue
		}

		// 查询公网 IP 归属地信息
		info, err := LookupIP(ipStr)
		if err != nil {
			logx.Info("查询失败")
			continue
		}

		fmt.Printf("IP: %s\n国家: %s\n省份: %s\n城市: %s\n运营商: %s\n经纬度: %s,%s\n",
			info.IP, info.Country, info.Region, info.City, info.ISP, info.Lat, info.Lng)
	}
}
