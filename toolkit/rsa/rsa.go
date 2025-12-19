package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
)

// RsaGenKey
//
//	 @Description: generates an RSA key pair and writes it to files.
//					bits: The length of the key in bits (e.g., 2048, 4096)
//	 @param bits
//	 @return error
func RsaGenKey(bits int) error {
	// 1. 生成私钥文件
	// GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	// 参数1: Reader是一个全局、共享的密码用强随机数生成器
	// 参数2: 秘钥的位数 - bit
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	// 2. MarshalPKCS1PrivateKey将rsa私钥序列化为ASN.1 PKCS#1 DER编码
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	// 3. Block代表PEM编码的结构, 对其进行设置
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	// 4. 创建文件
	privFile, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	// 5. 使用pem编码, 并将数据写入文件中
	err = pem.Encode(privFile, &block)
	if err != nil {
		return err
	}
	// 6. 最后的时候关闭文件
	defer privFile.Close()

	// 7. 生成公钥文件
	publicKey := privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err
	}
	block = pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	}
	pubFile, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	// 8. 编码公钥, 写入文件
	err = pem.Encode(pubFile, &block)
	if err != nil {
		panic(err)
		return err
	}
	defer pubFile.Close()

	return nil
}

// RSAEncrypt rsa加密
//
//	@Description:
//	@param data 要加密的数据
//	@param filename 公钥文件的路径
//	@return []byte
func RSAEncrypt(data, filename []byte) []byte {
	// 1. 根据文件名将文件内容从文件中读出
	file, err := os.Open(string(filename))
	if err != nil {
		return nil
	}
	// 2. 读文件
	info, _ := file.Stat()
	allText := make([]byte, info.Size())
	file.Read(allText)
	// 3. 关闭文件
	file.Close()

	// 4. 从数据中查找到下一个PEM格式的块
	block, _ := pem.Decode(allText)
	if block == nil {
		return nil
	}
	// 5. 解析一个DER编码的公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil
	}
	pubKey := pubInterface.(*rsa.PublicKey)

	// 6. 公钥加密
	result, _ := rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
	return result
}

// RSADecrypt rsa解密
//
//	@Description:
//	@param data 要解密的数据
//	@param filename 私钥文件的路径
//	@return []byte
func RSADecrypt(data, filename []byte) []byte {
	// 1. 根据文件名将文件内容从文件中读出
	file, err := os.Open(string(filename))
	if err != nil {
		return nil
	}
	// 2. 读文件
	info, _ := file.Stat()
	allText := make([]byte, info.Size())
	file.Read(allText)
	// 3. 关闭文件
	file.Close()
	// 4. 从数据中查找到下一个PEM格式的块
	block, _ := pem.Decode(allText)
	// 5. 解析一个pem格式的私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	// 6. 私钥解密
	result, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)

	return result
}

// B4RsaGenKey 生成 RSA 密钥对，并返回 base64 编码的 PEM 格式公钥和私钥字符串
func B4RsaGenKey(bits int) (base64PrivateKey, base64PublicKey string, err error) {
	// 1. 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	// 2. 编码私钥为 PKCS1 PEM 块
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	}
	privPEM := pem.EncodeToMemory(privBlock)

	// 3. 公钥编码为 PKIX 格式
	pubDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY", // 注意这里公钥格式是标准 PKIX（不是 RSA PUBLIC KEY）
		Bytes: pubDER,
	}
	pubPEM := pem.EncodeToMemory(pubBlock)

	// 4. 做 base64 编码
	base64PrivateKey = base64.StdEncoding.EncodeToString(privPEM)
	base64PublicKey = base64.StdEncoding.EncodeToString(pubPEM)

	return base64PrivateKey, base64PublicKey, nil
}

func B4RSAEncrypt(data []byte, base64PubKey string) []byte {
	// 1. 解 base64
	pemData, err := base64.StdEncoding.DecodeString(base64PubKey)
	if err != nil {
		return nil
	}

	// 2. 解 PEM
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil
	}

	// 3. 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil
	}
	pubKey := pubInterface.(*rsa.PublicKey)

	// 4. 加密
	result, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
	if err != nil {
		return nil
	}

	return result
}

func B4RSADecrypt(ciphertext []byte, base64PrivKey string) []byte {
	// 1. 解 base64
	pemData, err := base64.StdEncoding.DecodeString(base64PrivKey)
	if err != nil {
		return nil
	}

	// 2. 解 PEM
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil
	}

	// 3. 解析私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil
	}

	// 4. 解密
	result, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		return nil
	}

	return result
}
