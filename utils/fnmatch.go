package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"git.sxxfuture.net/taojiayi/super-cp/config"
	"git.sxxfuture.net/taojiayi/super-cp/utils/myglob"
)

/*
根据 options 规则，列出文件列表
*/
func MatchFiles(pattern interface{}) ([]string, error) {
	var matches []string
	var err error

	switch p := pattern.(type) {
	case config.SourcePattern:
		dotOption := p.Options["dot"]
		noCaseOption := p.Options["no-case"]
		matches, err = myglob.Match(os.DirFS("."), p.Glob, myglob.Options{
			Dot:    dotOption,
			NoCase: noCaseOption,
		})
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

/*
根据 rule 规则，判断某个路径是否匹配
*/
func MatchRuleFiles(pattern interface{}, targetPath string) bool {
	switch p := pattern.(type) {
	case config.RulePattern:
		matched := myglob.MatchPath(targetPath, p.Glob)
		return matched
	}

	return false
}
