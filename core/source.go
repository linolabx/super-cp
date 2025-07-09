package core

import (
	"mime"
	"os"
	"path/filepath"
	"strings"

	"git.sxxfuture.net/taojiayi/super-cp/utils"
	"github.com/gobwas/glob"
)

type Source struct {
	Pattern string `yaml:"pattern"`

	// strip prefix
	Strip string `yaml:"strip"`

	// keep unix hidden dirs and files, default not
	KeepDotFiles bool `yaml:"keep-dot-files"`
}

// walk dir and return all files match pattern and options
func (s *Source) WalkMatch() ([]*SourceFile, error) {
	matcher := glob.MustCompile(s.Pattern, '/')

	// get the max non-glob prefix to reduce walk range
	// foo/bar/baz**/?/* -> foo/bar/baz
	maxNonGlobPrefix := []string{}
	for _, part := range strings.Split(s.Pattern, "/") {
		if strings.Contains(part, "*") {
			break
		}
		if strings.Contains(part, "?") {
			break
		}
		maxNonGlobPrefix = append(maxNonGlobPrefix, part)
	}

	walkDir := filepath.Join(maxNonGlobPrefix...)
	utils.Verbosef("Walk Dir: %s", walkDir)

	files := []*SourceFile{}
	err := filepath.Walk(walkDir, func(localPath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		if !s.KeepDotFiles {
			for _, part := range strings.Split(localPath, "/") {
				if strings.HasPrefix(part, ".") {
					return nil
				}
			}
		}

		if !matcher.Match(localPath) {
			return nil
		}

		sFile := SourceFile{
			LocalPath:  localPath,
			RemotePath: localPath,
			Excluded:   false,
			Info:       f,
			Metadata:   map[string]string{},
		}

		if s.Strip != "" {
			sFile.RemotePath = strings.TrimPrefix(
				strings.TrimPrefix(sFile.RemotePath, s.Strip),
				"/",
			)
		}

		if mt := mime.TypeByExtension(filepath.Ext(sFile.LocalPath)); mt != "" {
			sFile.Metadata["Content-Type"] = mt
		}

		files = append(files, &sFile)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
