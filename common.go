package gosql

import (
	"github.com/tkdeng/regex"
)

// [^\w_\-]
func toAlphaNumeric(str string) string {
	return string(regex.Comp(`[^\w_\-]`).RepLit([]byte(str), []byte{}))
}

func sqlEscapeQuote(str string) string {
	return string(regex.Comp(`([\\"'\'])`).RepFunc([]byte(str), func(data func(int) []byte) []byte {
		return regex.JoinBytes(`\`, data(1))
	}))
}
