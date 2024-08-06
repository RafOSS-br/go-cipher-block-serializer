package goCipherBlockSerializer

import (
	"crypto/cipher"
	"reflect"
	"unsafe"
)

// Serialize serializes a cipher.Block to a byte slice.
func Serialize(cipher cipher.Block) []byte {
	val := reflect.ValueOf(cipher).Elem()
	size := val.Type().Size()
	ptr := unsafe.Pointer(val.UnsafeAddr())

	byteSlice := unsafe.Slice((*byte)(ptr), int(size))
	return byteSlice
}

// Deserialize deserializes a byte slice to a cipher.Block.
func Deserialize(data []byte, cipherType reflect.Type) cipher.Block {
	blockPtr := reflect.New(cipherType.Elem())
	val := blockPtr.Elem()
	size := val.Type().Size()
	ptr := unsafe.Pointer(val.UnsafeAddr())

	copy((*[1 << 30]byte)(ptr)[:size:size], data)

	return blockPtr.Interface().(cipher.Block)
}
