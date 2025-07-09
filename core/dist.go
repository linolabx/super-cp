package core

import (
	"os"

	"gopkg.in/yaml.v3"
)

// if dsn protocol equals to type, dist can be just a dsn string
type Dist struct {
	// 目标类型: s3, ssh, aliyun-oss
	Type string `yaml:"type"`
	// DSN to target
	DSN string `yaml:"dsn"`

	Uploader Uploader
}

func (d *Dist) init() error {
	d.DSN = os.ExpandEnv(d.DSN)
	uploader, err := GetUploader(d.Type, d.DSN)
	if err != nil {
		return err
	}
	d.Uploader = uploader

	return nil
}

func (d *Dist) UnmarshalYAML(value *yaml.Node) error {
	var dsn string
	if err := value.Decode(&dsn); err == nil {
		d.DSN = dsn
		return nil
	}

	var tmp Dist
	if err := value.Decode(&tmp); err != nil {
		return err
	}
	*d = tmp

	return nil
}
