package utils

import (
	"bytes"
	"encoding/json"
	"strings"
)

// Returns whenTrue if b is true; otherwise, returns whenFalse. IfElse should only be used when branches are either
// constant or precomputed as both branches will be evaluated regardless as to the value of b.
func IfElse[T any](b bool, whenTrue T, whenFalse T) T {
	if b {
		return whenTrue
	}
	return whenFalse
}

// Returns value if value is not the zero value of T; Otherwise, returns defaultValue. OrElse should only be used when
// defaultValue is constant or precomputed as its argument will be evaluated regardless as to the content of value.
func OrElse[T comparable](value T, defaultValue T) T {
	if value != *new(T) {
		return value
	}
	return defaultValue
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func StringifyJson(input any, prefix string, indent string) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent(prefix, indent)
	if _, ok := input.([]any); ok && len(input.([]any)) == 0 {
		return "[]", nil
	}
	if err := encoder.Encode(input); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}
