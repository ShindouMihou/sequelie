package sequelie

import (
	"bytes"
)

func ptr[T any](v T) *T {
	return &v
}

func insensitiveHasPrefix(s, prefix []byte) bool {
	return len(s) >= len(prefix) && bytes.EqualFold(s[0:len(prefix)], prefix)
}

func insensitiveHasSuffix(s, suffix []byte) bool {
	return len(s) >= len(suffix) && bytes.EqualFold(s[len(s)-len(suffix):], suffix)
}
