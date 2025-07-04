package source

import (
	"fmt"
	"strings"

	"git.sxxfuture.net/taojiayi/super-cp/config"
	"git.sxxfuture.net/taojiayi/super-cp/core"
	"git.sxxfuture.net/taojiayi/super-cp/utils"
)

/*
处理源文件，扫描文件，并返回文件列表
*/
func Scan(env config.Deployment) ([]core.SPFile, error) {
	// 源文件与目标文件（DSN）的映射
	files := make([]core.SPFile, 0)
	// 多个源，合并
	for _, src := range env.Source {
		// 使用 myglob 匹配文件, 返回文件列表
		sourceFiles, err := utils.MatchFiles(src.Pattern)
		if err != nil {
			fmt.Printf("match source pattern %s failed: %v\n", src.Pattern.Glob, err)
			continue
		}

		for _, sourceFile := range sourceFiles {
			if !utils.IsFile(sourceFile) {
				continue
			}

			// 应用路径前缀移除
			targetPath := sourceFile
			if src.Strip != "" {
				targetPath = strings.TrimPrefix(sourceFile, src.Strip)
				targetPath = strings.TrimPrefix(targetPath, "/")
			}

			files = append(files, *core.NewSPFile(sourceFile, targetPath, nil))

		}
	}

	return files, nil
}
