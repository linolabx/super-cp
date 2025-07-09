package core

import (
	"maps"

	"git.sxxfuture.net/taojiayi/super-cp/utils"
	"github.com/gobwas/glob"
)

// Rule 定义文件处理规则
type Rule struct {
	// glob
	Pattern string `yaml:"pattern"`

	// add headers to file
	Headers *map[string]string `yaml:"headers"`

	// remove file from upload list
	Exclude *bool `yaml:"exclude"`

	// set index of file
	Index *int `yaml:"index"`

	glob glob.Glob
}

func (r *Rule) init() error {
	r.glob = glob.MustCompile(r.Pattern, '/')

	if r.Headers != nil {
		for k, v := range *r.Headers {
			(*r.Headers)[k] = utils.Expand(v)
		}
	}

	return nil
}

func (r *Rule) Apply(file *SourceFile) error {
	if !r.glob.Match(file.RemotePath) {
		return nil
	}

	if r.Exclude != nil && *r.Exclude {
		return nil
	}

	if r.Headers != nil {
		maps.Copy(file.Metadata, *r.Headers)
	}

	if r.Index != nil {
		file.Index = *r.Index
	}

	return nil
}
