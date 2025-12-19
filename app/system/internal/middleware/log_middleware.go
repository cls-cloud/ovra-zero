package middleware

import (
	"net/http"
	"ovra/toolkit/ip"

	"github.com/zeromicro/go-zero/core/logx"
)

func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Info("请求前...")
		cuRip := ip.GetClientIP(r)
		logx.Info("请求 IP: ", cuRip)
		if ip.IsPrivateIP(cuRip) {
			logx.Info("请求 IP: 内网IP")
		}
		//ipData, _ := ip.LookupIP(cuRip)
		//fmt.Printf("IP: %s\n国家: %s\n省份: %s\n城市: %s\n运营商: %s\n经纬度: %s,%s\n",
		//	ipData.IP, ipData.Country, ipData.Region, ipData.City, ipData.ISP, ipData.Lat, ipData.Lng)
		next(w, r)
		logx.Info("请求后....")
	}
}
