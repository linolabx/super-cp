package utils

import (
	"github.com/gobwas/glob"
)

type GlobMatcher struct {
	pattern glob.Glob
}

func NewGlobMatcher(pattern string) *GlobMatcher {
	return &GlobMatcher{pattern: glob.MustCompile(pattern, '/')}
}

// it use gitignore pattern to match file
// so if it is 'excluded', it means matched
func (g *GlobMatcher) Match(path string) bool {
	return g.pattern.Match(path)
}
