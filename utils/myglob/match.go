package myglob

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// isHiddenPath 检查路径是否为隐藏路径
func isHiddenPath(path string, d fs.DirEntry) bool {
	// 检查当前文件/目录名称
	if strings.HasPrefix(d.Name(), ".") && d.Name() != "." && d.Name() != ".." {
		return true
	}
	// 检查路径的每一部分
	parts := strings.Split(filepath.ToSlash(path), "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") && part != "." && part != ".." {
			return true
		}
	}
	return false
}

/*
根据 pattern + options 规则，列出文件列表（包括文件和目录）
*/
func Match(fsys fs.FS, pattern string, opts Options) ([]string, error) {
	var matches []string

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// 跳过隐藏文件/目录
		if !opts.Dot && isHiddenPath(path, d) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		target := path
		matchPattern := pattern
		if opts.NoCase {
			matchPattern = strings.ToLower(pattern)
			target = strings.ToLower(path)
		}

		ok, err := matchRegex(matchPattern, target)
		if err != nil {
			return nil
		}
		if ok {
			matches = append(matches, path)
		}
		return nil
	})

	return matches, err
}

/*
判断某个路径是否匹配 pattern 规则
*/
func MatchPath(path string, pattern string) bool {
	matchPattern := pattern
	ok, err := matchRegex(matchPattern, path)
	if err != nil {
		return false
	}
	return ok
}
