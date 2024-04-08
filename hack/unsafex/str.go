//go:build !go1.20
// +build !go1.20

package unsafex

import (
	"unsafe"
)

// String2Bytes converts string to []byte without memory allocation
// Warning: read-only
func StringToSlice(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// Bytes2String converts []byte to string without memory allocation
func SliceToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
