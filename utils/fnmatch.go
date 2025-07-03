package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"git.sxxfuture.net/taojiayi/super-cp/config"
	"github.com/bmatcuk/doublestar/v4"
)

// 列出 source 文件列表
func MatchSourceFiles(pattern interface{}) ([]string, error) {
	var matches []string
	var err error

	switch p := pattern.(type) {
	case string: // 普通模式（不会走这里了，因为 unmarshal 的时候已经处理了）
		matches, err = doublestar.Glob(os.DirFS("."), p)
	case config.SourcePattern: // glob 模式
		matches, err = doublestar.Glob(os.DirFS("."), p.Glob)
	default:
		return nil, fmt.Errorf("unsupported pattern type: %T", pattern)
	}

	if err != nil {
		return nil, err
	}

	var files []string
	for _, rel := range matches {
		fullPath := filepath.Join(".", rel)
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
			files = append(files, fullPath)
		}
	}
	return files, nil
}

func MatchRuleFiles(pattern interface{}, targetPath string) bool {
	switch p := pattern.(type) {
	case string: // 普通模式（不会走这里了，因为 unmarshal 的时候已经处理了）
		matched, err := doublestar.PathMatch(p, targetPath)
		if err != nil {
			fmt.Printf("解析规则模式 %s 失败: %v\n", p, err)
			return false
		}
		return matched
	case config.RulePattern: // glob 模式
		matched, err := doublestar.PathMatch(p.Glob, targetPath)
		if err != nil {
			fmt.Printf("解析规则模式 %s 失败: %v\n", p.Glob, err)
			return false
		}
		return matched
	}

	return false
}
