package util

import (
	"encoding/binary"
	"unicode/utf8"
)

func IntToBytes(num int) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(num))
	return b
}

func ByteToInt(b []byte) int {
	return int(binary.BigEndian.Uint32(b))
}

func Utf8Len(s string) int {
	return utf8.RuneCountInString(s)
}
