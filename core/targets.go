package core

import (
	"fmt"
	"net/url"
)

type Uploader interface {
	Upload(files []*SourceFile) error
}

type UploaderFactory func(dsn string) (Uploader, error)

var uploaders = map[string]UploaderFactory{}

func RegisterUploader(name string, factory UploaderFactory) {
	uploaders[name] = factory
}

func GetUploader(name, dsn string) (Uploader, error) {
	if name != "" {
		if factory, ok := uploaders[name]; ok {
			return factory(dsn)
		}
	}

	dsnUrl, err := url.Parse(dsn)
	if err == nil {
		if factory, ok := uploaders[dsnUrl.Scheme]; ok {
			return factory(dsn)
		}
	}

	for _, factory := range uploaders {
		if uploader, err := factory(dsn); err == nil {
			return uploader, nil
		}
	}

	return nil, fmt.Errorf("uploader not found: %s", dsn)
}
