package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// GenerateKey 生成ed25519密钥对
func GenerateKey() (publicKeyStr string, privateKeyStr string, err error) {
	// 生成密钥对
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	// 编码公钥
	publicBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}
	publicPemFile := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicBytes,
	})

	// 编码私钥
	privateBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}
	privatePemFile := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateBytes,
	})

	publicKeyStr = string(publicPemFile)
	privateKeyStr = string(privatePemFile)

	return
}

// ParsePublicKey 解码公钥
func ParsePublicKey(data string) (publicKey ed25519.PublicKey, err error) {
	// pem解密
	block, _ := pem.Decode([]byte(data))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("解码公钥时无法解析 PEM 数据: %s", data)
	}
	// x509解密
	ed25519Key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("公钥解码错误: %v", err)
	}

	publicKey, ok := ed25519Key.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("无效的 ed25519 公钥: %s", data)
	}
	return
}

// ParsePrivateKey 解码秘钥
func ParsePrivateKey(data string) (privateKey ed25519.PrivateKey, err error) {
	// pem解码
	block, _ := pem.Decode([]byte(data))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("解码秘钥时无法解析 PEM 数据: %s", data)
	}
	// x509解码
	ed25519Key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("私钥解码错误: %v", err)
	}

	privateKey, ok := ed25519Key.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("无效的 ed25519 私钥: %s", data)
	}
	return
}

// ED25519Signer ED25519签名工具
type ED25519Signer struct {
	privateKey ed25519.PrivateKey
}

// Sign 签名数据
func (signer *ED25519Signer) Sign(data []byte) (signatureB64 string, err error) {
	signature := ed25519.Sign(signer.privateKey, data)
	signatureB64 = base64.StdEncoding.EncodeToString(signature)
	return
}

// ED25519Verifier ED25519校验工具
type ED25519Verifier struct {
	publicKey ed25519.PublicKey
}

// Verify 验证签名
func (verifier *ED25519Verifier) Verify(data []byte, signatureB64 string) (valid bool, err error) {
	// 解码 Base64 编码的签名
	signature, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return false, fmt.Errorf("使用 ED2519 校验数据数据时, Base64 解码错误: %v", err)
	}
	valid = ed25519.Verify(verifier.publicKey, data, signature)
	return
}
