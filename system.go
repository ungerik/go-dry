package dry

import (
	"unsafe"
)

func SystemIsLittleEndian() bool {
	var word uint16 = 1
	littlePtr := (*uint8)(unsafe.Pointer(&word))
	return (*littlePtr) == 1
}

func SystemIsBigEndian() bool {
	return !SystemIsLittleEndian()
}
