package utils

import (
	"os"
)

// IsFile 检查是否为普通文件
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}
