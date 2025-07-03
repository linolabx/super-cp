package utils

import (
	"os"
	"regexp"
	"strings"
)

/*
替换环境变量
*/
func ReplaceEnvVars(s string) string {
	env := os.Environ()
	envMap := make(map[string]string)

	// 将环境变量转换为 map 形式方便查找
	for _, v := range env {
		split := strings.SplitN(v, "=", 2)
		if len(split) >= 2 {
			envMap[split[0]] = split[1]
		}
	}

	re := regexp.MustCompile(`\$([a-zA-Z_]\w*)`)
	matches := re.FindAllStringSubmatch(s, -1)
	for _, match := range matches {
		if envValue, exists := envMap[match[1]]; exists {
			s = strings.ReplaceAll(s, match[0], envValue)
		}
	}
	return s
}

// IsFile 检查是否为普通文件
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}
