package core

type Target interface {
	Init(dsn string) error
	Upload(transfer SPFile) error
}

var Targets = map[string]Target{}
