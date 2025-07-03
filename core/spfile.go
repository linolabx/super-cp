package core

type SPFile struct {
	LocalPath  string
	RemotePath string
	Metadata   map[string]string
}

func NewSPFile(localPath string, remotePath string, metadata map[string]string) *SPFile {
	return &SPFile{
		LocalPath:  localPath,
		RemotePath: remotePath,
		Metadata:   metadata,
	}
}
