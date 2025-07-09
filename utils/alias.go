package utils

import (
	"strings"
)

var StaticAliases = map[string]string{
	"@forever": "public, max-age=31536000, immutable",
	"@year":    "public, max-age=31536000",
	"@month":   "public, max-age=2592000",
	"@week":    "public, max-age=604800",
	"@day":     "public, max-age=86400",
	"@hour":    "public, max-age=3600",
	"@minute":  "public, max-age=60",
	"@second":  "public, max-age=1",
}

func Expand(s string) string {
	for k, v := range StaticAliases {
		s = strings.ReplaceAll(s, k, v)
	}

	return s
}
