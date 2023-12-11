package rsa

import (
	"github.com/agclqq/goencryption"
	"os"
)

// RsaEncrypt
//
//	@Description: Rsa加密
//	@param text	需要加密的字符串
//	@return encryptStr	加密后的密文字符串
//	@return err	错误信息
func RsaEncrypt(text string) (encryptStr string, err error) {
	// 读取公钥
	pubKey, err1 := os.ReadFile("internal/service/rsa/pub.key")
	if err1 != nil {
		//app.LLog.Info("读取公钥失败:", err1)
		return encryptStr, err1
	}
	//使用公钥加密
	encryptByte, err2 := goencryption.PubKeyEncrypt(pubKey, []byte(text))
	if err2 != nil {
		//app.LLog.Info("加密失败:", err2)
		return encryptStr, err2
	}
	// 密文字符串转Base64
	return goencryption.Base64Encode(encryptByte), nil
}

// RsaDecrypt
//
//	@Description: RSA解密
//	@param decryptStr	rsa加密过的密文
//	@return text	解密后的字符串
//	@return err	错误信息
func RsaDecrypt(decryptStr string) (text string, err error) {
	// 读取私钥
	priKey, err1 := os.ReadFile("internal/service/rsa/pri.key")
	if err1 != nil {
		//app.LLog.Info("读取密钥失败:", err1)
		return text, err1
	}
	// 密文字符串转Base64
	cipherText, err2 := goencryption.Base64Decode(decryptStr)
	if err2 != nil {
		//app.LLog.Info("Base64Decode失败:", err2)
		return text, err2
	}
	//使用私钥解密
	textByte, err3 := goencryption.PrvKeyDecrypt(priKey, cipherText)
	if err3 != nil {
		//app.LLog.Info("解密失败:", err3)
		return text, err3
	}
	return string(textByte), nil
}
