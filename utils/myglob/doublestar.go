package myglob

import (
	"path/filepath"
	"strings"
)

/*
判断某个路径是否匹配 pattern 规则
*/
func matchRegex(pattern, path string) (bool, error) {
	patternParts := strings.Split(pattern, string(filepath.Separator))
	pathParts := strings.Split(path, string(filepath.Separator))
	return matchParts(patternParts, pathParts), nil
}

func matchParts(patParts, pathParts []string) bool {
	if len(patParts) == 0 {
		return len(pathParts) == 0
	}
	if patParts[0] == "**" {
		for i := 0; i <= len(pathParts); i++ {
			if matchParts(patParts[1:], pathParts[i:]) {
				return true
			}
		}
		return false
	}
	if len(pathParts) == 0 {
		return false
	}
	ok, _ := filepath.Match(patParts[0], pathParts[0])
	if !ok {
		return false
	}
	return matchParts(patParts[1:], pathParts[1:])
}
