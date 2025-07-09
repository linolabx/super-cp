package utils

import (
	"fmt"
	"testing"
)

func TestGlobMatcher(t *testing.T) {
	fmt.Println(NewGlobMatcher("*/*").Match("dist/index.html"))
}
