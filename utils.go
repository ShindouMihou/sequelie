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
