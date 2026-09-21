// Package unsafex provides zero-copy conversions between string and []byte.
//
// It relies on unsafe.String/unsafe.StringData/unsafe.Slice/unsafe.SliceData,
// so it requires go1.20 or later and the results must be treated as read-only.
package unsafex

import (
	"unsafe"
)

// SliceToString converts []byte to string without data copy.
// Warning: the returned string shares the memory of b, writing to b afterwards
// breaks the immutability of strings.
func SliceToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// StringToSlice converts string to []byte without data copy.
// Warning: the returned slice must not be written to.
func StringToSlice(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
