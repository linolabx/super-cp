package core

import "os"

type SourceFile struct {
	LocalPath  string
	RemotePath string
	Excluded   bool
	Info       os.FileInfo
	Index      int
	Metadata   map[string]string
}
