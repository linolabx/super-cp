package utils

import (
	"testing"
)

func TestIsFile(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{path: "helper_test.go", want: true},
		{path: "helper.go", want: true},
		{path: "helper_test.go.txt", want: false},
	}
	for _, test := range tests {
		if got := IsFile(test.path); got != test.want {
			t.Errorf("IsFile(%q) = %v, want %v", test.path, got, test.want)
		}
	}
}
