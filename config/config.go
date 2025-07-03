package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Deployment 定义传输策略
type Deployment struct {
	Source []Source `yaml:"source"`
	Dist   Dist     `yaml:"dist"`
	Rules  []Rule   `yaml:"rules"`
}

// Source 定义要传输的文件
type Source struct {
	Pattern SourcePattern `yaml:"pattern"` // 可以是 string 或 SourcePattern（glob 模式）
	Strip   string        `yaml:"strip"`
}

type SourcePattern struct {
	Glob    string         `yaml:"glob"`    // 路径，例如 "**/*.html"，"**/*.js"
	Options map[string]any `yaml:"options"` // 例如 dot: true，noext: true 等
}

// Dist 定义目标存储
type Dist struct {
	Type string `yaml:"type"`
	DSN  string `yaml:"dsn"`
}

// Rule 定义文件处理规则
type Rule struct {
	Pattern RulePattern `yaml:"pattern"` // 可以是 string 或 RulePattern（glob 模式）

	Headers      map[string]string `yaml:"headers"`        // 规则头
	AutoMimeType bool              `yaml:"auto-mime-type"` // 是否自动根据文件名设置 Content-Type

	Exclude bool `yaml:"exclude"` // 是否排除上传
}

type RulePattern struct {
	Glob    string         `yaml:"glob"`    // 路径，例如 "**/*.html"，"**/*.js"
	Options map[string]any `yaml:"options"` // 例如 dot: true，noext: true 等
}

func LoadConfig(path string) (Config, error) {
	// 读取并解析配置
	configData, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("read config file failed: %v\n", err)
		os.Exit(1)
	}

	config := Config{}
	if err := yaml.Unmarshal(configData, &config); err != nil {
		fmt.Printf("parse config file failed: %v\n", err)
		os.Exit(1)
	}

	return config, nil
}

type Config struct {
	Environments map[string]Deployment `yaml:"environments"`
}

func (c *Config) GetEnv(name string) (bool, Deployment) {
	env, ok := c.Environments[name]
	return ok, env
}

func (sp *SourcePattern) UnmarshalYAML(value *yaml.Node) error {
	// Try to unmarshal as string
	var glob string
	if err := value.Decode(&glob); err == nil {
		sp.Glob = glob
		sp.Options = nil
		return nil
	}
	// Try to unmarshal as struct
	type alias SourcePattern // prevent recursion
	var tmp alias
	if err := value.Decode(&tmp); err != nil {
		return err
	}
	*sp = SourcePattern(tmp)
	return nil
}

func (rp *RulePattern) UnmarshalYAML(value *yaml.Node) error {
	// Try to unmarshal as string
	var glob string
	if err := value.Decode(&glob); err == nil {
		rp.Glob = glob
		rp.Options = nil
		return nil
	}
	// Try to unmarshal as struct
	type alias RulePattern // prevent recursion
	var tmp alias
	if err := value.Decode(&tmp); err != nil {
		return err
	}
	*rp = RulePattern(tmp)
	return nil
}
