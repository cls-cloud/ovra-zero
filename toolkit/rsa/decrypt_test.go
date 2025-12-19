package rsa

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"testing"
)

var (
	// 请求体：使用 AES CBC 加密 + Base64 编码，粘贴密文
	mockEncryptedKey = `bcmhSKvs5e0kN/Rrsk/ncyMnhTsYI5iGjdt/B0lIAFpst+gTYk2eybxy0IGJ9Yms+Lglx+ASg9cplGwmKGT/Ng==`
	// 前端的body
	mockEncryptedBody = `pSY5hs/PufC/jPlFx01hHlJpz8+VKeq6ve9LGtsq41AenbyawKPhMKbuQFt8+N/4cEYI838ixVb9HkHpFkO7+6OxIQrT5yHLhJJ5VCQXCoX9fu8NXcmeBBpktd9bL2ui5xIIOlE//dDcCe/AW7KHby53iSNSHV27HAVn7IHgpHGbdMz1Sa8ecoQd0b/rmyW91PmovIy6w6PwADfB3BbWRHNtU6zPZ8OJyJkI8GyNLSg=`
)

// 模拟前端使用的 AES 密钥 Base64 解码
func decryptRSA(encryptedBase64 string, privateKey *rsa.PrivateKey) ([]byte, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(nil, privateKey, encryptedBytes)
}

func decryptAESCBCBase64(cipherTextBase64 string, aesKey []byte, iv []byte) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	if len(cipherText)%aes.BlockSize != 0 {
		return "", errors.New("cipherText is not a multiple of block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	mode.CryptBlocks(plainText, cipherText)
	plainText = pkcs7Unpadding(plainText)
	return string(plainText), nil
}

func pkcs7Unpadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

func loadRSAPrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func loadRSAPrivateKeyFromBase64(base64Key string) (*rsa.PrivateKey, error) {
	derBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, fmt.Errorf("base64 decode failed: %w", err)
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(derBytes)
	if err != nil {
		return nil, fmt.Errorf("parse PKCS#8 private key failed: %w", err)
	}
	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}
	return rsaKey, nil
}

func TestDecrypt(t *testing.T) {
	// 1. 加载私钥
	privateKeyBase64 := `MIIBVAIBADANBgkqhkiG9w0BAQEFAASCAT4wggE6AgEAAkEAmc3CuPiGL/LcIIm7zryCEIbl1SPzBkr75E2VMtxegyZ1lYRD+7TZGAPkvIsBcaMs6Nsy0L78n2qh+lIZMpLH8wIDAQABAkEAk82Mhz0tlv6IVCyIcw/s3f0E+WLmtPFyR9/WtV3Y5aaejUkU60JpX4m5xNR2VaqOLTZAYjW8Wy0aXr3zYIhhQQIhAMfqR9oFdYw1J9SsNc+CrhugAvKTi0+BF6VoL6psWhvbAiEAxPPNTmrkmrXwdm/pQQu3UOQmc2vCZ5tiKpW10CgJi8kCIFGkL6utxw93Ncj4exE/gPLvKcT+1Emnoox+O9kRXss5AiAMtYLJDaLEzPrAWcZeeSgSIzbL+ecokmFKSDDcRske6QIgSMkHedwND1olF8vlKsJUGK3BcdtM8w4Xq7BpSBwsloE=`

	// 解析成 RSA 私钥对象
	privateKey, err := loadRSAPrivateKeyFromBase64(privateKeyBase64)
	if err != nil {
		t.Fatal("解析私钥失败:", err)
	}
	if err != nil {
		panic("加载私钥失败：" + err.Error())
	}

	// 2. 解密 encrypt-key，获取 AES 密钥（注意里面是 Base64 编码的 key）
	aesKeyBase64, err := decryptRSA(mockEncryptedKey, privateKey)
	if err != nil {
		panic("RSA 解密失败：" + err.Error())
	}

	aesKey, err := base64.StdEncoding.DecodeString(string(aesKeyBase64))
	if err != nil {
		panic("Base64 解码 AES key 失败：" + err.Error())
	}

	fmt.Printf("获取 AES Key: %s\n", string(aesKey))

	// 3. 解密请求体
	iv := []byte("A-16-Bytes-String") // 默认 CryptoJS 中固定 IV，需确认
	decryptedText, err := decryptAESCBCBase64(mockEncryptedBody, aesKey, iv)
	if err != nil {
		panic("AES 解密失败：" + err.Error())
	}

	fmt.Println("解密后的请求体内容:")
	fmt.Println(decryptedText)
}
