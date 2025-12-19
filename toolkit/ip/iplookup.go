package ip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

// IP查询通过阿里云市场购买获取，当前apiKey仅为测试使用
const (
	ipAPIHost = "https://c2ba.api.huachen.cn"
	ipAPIPath = "/ip"
	ipAppCode = "d6cb28dfab5c4e8fb521144833327251" // 替换为真实 AppCode
)

type IPGeoData struct {
	IP         string `json:"ip"`
	LongIP     string `json:"long_ip"`
	ISP        string `json:"isp"`
	Area       string `json:"area"`
	RegionID   string `json:"region_id"`
	Region     string `json:"region"`
	CityID     string `json:"city_id"`
	City       string `json:"city"`
	District   string `json:"district"`
	DistrictID string `json:"district_id"`
	CountryID  string `json:"country_id"`
	Country    string `json:"country"`
	Lat        string `json:"lat"`
	Lng        string `json:"lng"`
}

type IPGeoResp struct {
	Ret   int       `json:"ret"`
	Msg   string    `json:"msg"`
	Data  IPGeoData `json:"data"`
	LogID string    `json:"log_id"`
}

func LookupIP(ip string) (*IPGeoData, error) {
	query := url.Values{}
	query.Set("ip", ip)
	fullURL := fmt.Sprintf("%s%s?%s", ipAPIHost, ipAPIPath, query.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request failed: %w", err)
	}
	req.Header.Set("Authorization", "APPCODE "+ipAppCode)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result IPGeoResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}

	if result.Ret != 200 {
		return nil, fmt.Errorf("API error: %s", result.Msg)
	}

	return &result.Data, nil
}

func IsPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	if ip.IsLoopback() {
		return true
	}

	privateBlocks := []*net.IPNet{
		// IPv4 私有地址段
		{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
		{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
		// IPv6 ULA（本地地址 fc00::/7）
		{IP: net.ParseIP("fc00::"), Mask: net.CIDRMask(7, 128)},
		// IPv6 链路本地地址（fe80::/10）
		{IP: net.ParseIP("fe80::"), Mask: net.CIDRMask(10, 128)},
	}

	for _, block := range privateBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}
