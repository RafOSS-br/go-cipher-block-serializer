package goCipherBlockSerializer

import (
	"crypto/cipher"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"unsafe"
)

type Block struct {
	Enc []uint32
	Dec []uint32
}

// FromJSON reads a JSON and returns a Block.
func NewBlockFromJson(reader io.ByteReader) (Block, error) {
	var b Block
	var buffer []byte
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return Block{}, err
		}
		buffer = append(buffer, b)
	}
	err := json.Unmarshal(buffer, &b)
	if err != nil {
		return Block{}, err
	}
	return b, nil
}

// Json writes the JSON representation of a Block.
func (b Block) Json(f io.Writer) error {
	bytes, err := json.Marshal(b)
	if err != nil {
		return err
	}
	f.Write(bytes)
	return nil
}

// Serialize serializes a cipher.Block to a byte slice.
func Serialize(cipherBlock cipher.Block) (Block, error) {
	if cipherBlock == nil {
		return Block{}, errors.New("cipher block is nil")
	}

	val := reflect.ValueOf(cipherBlock).Elem()

	// Traverse nested structs until we find the actual enc/dec fields
	for val.Kind() == reflect.Struct && val.NumField() == 1 {
		val = val.FieldByIndex([]int{0})
	}

	encValue := val.FieldByName("enc")
	if !encValue.IsValid() || encValue.Kind() != reflect.Slice || encValue.IsNil() {
		return Block{}, errors.New("enc field is invalid or nil")
	}

	decValue := val.FieldByName("dec")
	if !decValue.IsValid() || decValue.Kind() != reflect.Slice || decValue.IsNil() {
		return Block{}, errors.New("dec field is invalid or nil")
	}

	enc := *(*[]uint32)(unsafe.Pointer(encValue.UnsafeAddr()))
	dec := *(*[]uint32)(unsafe.Pointer(decValue.UnsafeAddr()))

	return Block{Enc: enc, Dec: dec}, nil
}

// Deserialize deserializes a byte slice to a cipher.Block.
func Deserialize(block Block, cipherType reflect.Type) (cipher.Block, error) {
	if cipherType == nil {
		return nil, errors.New("cipherType is nil")
	}

	if cipherType.Kind() == reflect.Pointer {
		cipherType = cipherType.Elem()
	}

	if cipherType.Kind() != reflect.Struct {
		return nil, errors.New("cipherType is not a struct")
	}

	// Create a new instance of the cipher.Block using the provided type
	newBlockPtr := reflect.NewAt(cipherType, unsafe.Pointer(&block))

	// Convert the newBlockPtr to a cipher.Block interface
	cipherBlock, ok := newBlockPtr.Interface().(cipher.Block)
	if !ok {
		return nil, errors.New("failed to convert to cipher.Block")
	}

	return cipherBlock, nil
}
