package crypto

// 加密工具
type Encoder interface {
	Encrypt(pt []byte) (ctB64 string, err error)
}

// 解密工具
type Decoder interface {
	Decrypt(ctB64 string) (pt []byte, err error)
}

// 签名工具
type Signer interface {
    Sign(data []byte) (signatureB64 string, err error)
}

// 校验工具
type Verifier interface {
    Verify(data []byte, signatureB64 string) (valid bool, err error)
}
