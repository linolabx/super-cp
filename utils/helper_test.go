package utils

import (
	"os"
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

func TestReplaceEnvVars(t *testing.T) {
	os.Setenv("S3_DSN", "s3://accessKey:secretKey@endpoint/bucket")
	tests := []struct {
		s    string
		want string
	}{
		{s: "$S3_DSN", want: "s3://accessKey:secretKey@endpoint/bucket"},
		{s: "S3_DSN_NOT_EXIST", want: "S3_DSN_NOT_EXIST"},
	}
	for _, test := range tests {
		if got := ReplaceEnvVars(test.s); got != test.want {
			t.Errorf("ReplaceEnvVars(%q) = %v, want %v", test.s, got, test.want)
		}
	}
}
