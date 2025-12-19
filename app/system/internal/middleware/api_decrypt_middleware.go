package middleware

import (
	"fmt"
	"net/http"
	"ovra/app/system/internal/config"
	"ovra/toolkit/rsa"
)

type ApiDecryptMiddleware struct {
	c config.Config
}

func NewApiDecryptMiddleware(c config.Config) *ApiDecryptMiddleware {
	return &ApiDecryptMiddleware{c: c}
}

func (m *ApiDecryptMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if m.c.ApiDecrypt.Enabled {

			// 加密前数据
			data := []byte("Hello RSA")
			// 假设 base64PubKey 是你存储在配置或环境变量中的 PEM 格式公钥的 base64 编码字符串
			encrypted := rsa.B4RSAEncrypt(data, m.c.ApiDecrypt.PublicKey)
			fmt.Println("加密后的数据:", encrypted)
			// 解密
			decrypted := rsa.B4RSADecrypt(encrypted, m.c.ApiDecrypt.PrivateKey)
			fmt.Println("解密后的数据", string(decrypted))
			//请求参数解密

			//响应参数加密
		}
		next(w, r)
	}
}
