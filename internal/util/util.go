package util

import (
	"reflect"
	"unsafe"
)

// Global variable for size
var ptrSize = func() uintptr {
	if unsafe.Sizeof(uintptr(0)) == 4 {
		return 4 // 32-bit system
	}
	return 8 // 64-bit system
}()

func Zero[T any](obj *T) {
	size := reflect.TypeOf(obj).Elem().Size()
	ptr := unsafe.Pointer(obj)

	{
		var i uintptr
		// Perform XOR swap byte by byte
		for ; i+ptrSize <= size; i += ptrSize {
			*(*uintptr)(ptr) = 0

			ptr = unsafe.Add(ptr, int(ptrSize))
		}

		for ; i < size; i++ {
			*(*byte)(ptr) = 0

			ptr = unsafe.Add(ptr, 1)
		}
	}
}
