package file

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CalculateSHA256 计算文件的 SHA-256 哈希值
func CalculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Copy 复制文件
func Copy(src, dest string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	destDir := filepath.Dir(dest)
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		return os.MkdirAll(destDir, 0755)
	}
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 复制文件内容
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}
	return nil
}

// CopyAndVerify 复制文件并校验
func CopyAndVerify(src, dest, srcHash string) error {
	err := Copy(src, dest)
	if err != nil {
		return err
	}

	destHash, err := CalculateSHA256(dest)
	if err != nil {
		return err
	}
	// 比较哈希值
	if srcHash != destHash {
		return fmt.Errorf("hash mismatch: source and destination files are not identical")
	}
	return nil
}

// CopyFS 从嵌入文件中复
func CopyFS(fs embed.FS, fsName string, dest string) error {
	dir := filepath.Dir(dest)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	// 从嵌入的文件系统中提取文件
	data, err := fs.ReadFile(fsName)
	if err != nil {
		return err
	}

	// 将文件写入目标路径
	err = os.WriteFile(dest, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
