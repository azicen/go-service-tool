package ecc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

// GenerateKey 生成ECC密钥对
func GenerateKey() (publicKeyStr string, privateKeyStr string, err error) {
	// 生成秘钥
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}
	// x509编码
	eccPrivateKey, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}
	// pem编码
	privateBlock := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: eccPrivateKey,
	}

	// 获取公钥
	publicKey := privateKey.PublicKey
	// x509编码
	eccPublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return "", "", err
	}
	// pem编码
	block := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: eccPublicKey,
	}
	publicKeyStr = string(pem.EncodeToMemory(&block))
	privateKeyStr = string(pem.EncodeToMemory(&privateBlock))

	return
}

// ParsePrivateKey 解码秘钥
func ParsePrivateKey(data string) (privateKey *ecdsa.PrivateKey, err error) {
	// pem解码
	block, _ := pem.Decode([]byte(data))
	// x509解码
	privateKey, err = x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return
}

// ParsePublicKey 解码公钥
func ParsePublicKey(data string) (publicKey *ecdsa.PublicKey, err error) {
	// pem解密
	block, _ := pem.Decode([]byte(data))
	// x509解密
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey = publicInterface.(*ecdsa.PublicKey)
	return
}

// ECCEncoder ECC加密工具
type ECCEncoder struct {
	publicKey *ecdsa.PublicKey
}

// Encrypt ECC加密
func (encoder *ECCEncoder) Encrypt(pt []byte) (ctB64 string, err error) {
	// ecdsa key 转换 ecies key
	eciesPublicKey := ecies.ImportECDSAPublic(encoder.publicKey)
	// 加密
	ct, err := ecies.Encrypt(rand.Reader, eciesPublicKey, pt, nil, nil)
	if err != nil {
		return "", err
	}
	// base64编码
	ctB64 = base64.StdEncoding.EncodeToString(ct)
	return
}

// ECCDecoder ECC解密工具
type ECCDecoder struct {
	privateKey *ecdsa.PrivateKey
}

// Decrypt ECC解密
func (decoder *ECCDecoder) Decrypt(ctB64 string) (pt []byte, err error) {
	// base64解码
	ct, err := base64.StdEncoding.DecodeString(ctB64)
	if err != nil {
		return nil, err
	}
	// ecdsa key 转换 ecies key
	eciesPrivateKey := ecies.ImportECDSA(decoder.privateKey)
	// 解密
	pt, err = eciesPrivateKey.Decrypt(ct, nil, nil)
	return
}
