package byteutil

import "unsafe"

func SubBytes(bytes []byte, start, end int) []byte {
	if len(bytes) == 0 {
		return []byte{}
	}
	if end >= len(bytes) {
		end = len(bytes) - 1
	}
	return bytes[start:end]
}

func SliceByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToSliceByte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}