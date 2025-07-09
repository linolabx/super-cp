package core

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Job struct {
	Source *Source `yaml:"source"`
	Dist   *Dist   `yaml:"dist"`
	Rules  []*Rule `yaml:"rules"`
}

type Config struct {
	Jobs map[string]*Job `yaml:"jobs"`
}

func MustLoadConfig(path string) Config {
	configData, err := os.ReadFile(path)
	if err != nil {
		log.Panicf("failed to read config file %s: %v", path, err)
	}

	config := Config{}
	if err := yaml.Unmarshal(configData, &config); err != nil {
		log.Panicf("failed to parse config file %s: %v", path, err)
	}

	for jobName, job := range config.Jobs {
		if err := job.Dist.init(); err != nil {
			log.Panicf("failed to init dist for job %s: %v", jobName, err)
		}
		for _, rule := range job.Rules {
			if err := rule.init(); err != nil {
				log.Panicf("failed to init rule for job %s: %s: %v", jobName, rule.Pattern, err)
			}
		}
	}

	return config
}
