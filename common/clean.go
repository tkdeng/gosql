package common

import (
	"bytes"
	"strings"
)

// Clean.Str will sanitizes a string to valid UTF-8
func Clean[T interface{ string | []byte }](val T) T {
	valT := ToInterface{val}.Val

	if v, ok := valT.(string); ok {
		return ToInterface{strings.ToValidUTF8(v, "")}.Val.(T)
	} else if v, ok := valT.([]byte); ok {
		return ToInterface{bytes.ToValidUTF8(v, []byte{})}.Val.(T)
	}

	return NullType[T]{}.Null
}
