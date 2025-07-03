package rules

import (
	"mime"
	"path/filepath"

	"git.sxxfuture.net/taojiayi/super-cp/config"
	"git.sxxfuture.net/taojiayi/super-cp/core"
	"git.sxxfuture.net/taojiayi/super-cp/utils"
)

/*
处理规则，根据规则处理 spFile 文件列表

返回处理后的文件列表
*/
func ProcessRule(rules []config.Rule, spFiles []core.SPFile) ([]core.SPFile, error) {
	var processed []core.SPFile

	for _, file := range spFiles {
		var matched bool
		metadata := map[string]string{}

		for _, rule := range rules {
			if !utils.MatchRuleFiles(rule.Pattern, file.RemotePath) {
				continue
			}

			if rule.Exclude {
				matched = false
				break // 被排除，不处理该文件
			}

			matched = true
			for k, v := range rule.Headers {
				metadata[k] = v
			}

			if rule.AutoMimeType {
				ext := filepath.Ext(file.RemotePath)
				if ct := mime.TypeByExtension(ext); ct != "" {
					metadata["Content-Type"] = ct
				}
			}
		}

		if matched {
			file.Metadata = metadata
			processed = append(processed, file)
		}
	}

	return processed, nil
}
