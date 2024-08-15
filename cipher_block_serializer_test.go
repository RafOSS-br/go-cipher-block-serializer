package goCipherBlockSerializer

import (
	"bytes"
	"crypto/aes"
	"io"
	"reflect"
	"testing"
)

// Common values for tests.

var commonInput = []byte{
	0x6b, 0xc1, 0xbe, 0xe2, 0x2e, 0x40, 0x9f, 0x96, 0xe9, 0x3d, 0x7e, 0x11, 0x73, 0x93, 0x17, 0x2a,
	0xae, 0x2d, 0x8a, 0x57, 0x1e, 0x03, 0xac, 0x9c, 0x9e, 0xb7, 0x6f, 0xac, 0x45, 0xaf, 0x8e, 0x51,
	0x30, 0xc8, 0x1c, 0x46, 0xa3, 0x5c, 0xe4, 0x11, 0xe5, 0xfb, 0xc1, 0x19, 0x1a, 0x0a, 0x52, 0xef,
	0xf6, 0x9f, 0x24, 0x45, 0xdf, 0x4f, 0x9b, 0x17, 0xad, 0x2b, 0x41, 0x7b, 0xe6, 0x6c, 0x37, 0x10,
}

var commonKey128 = []byte{0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6, 0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c}

var commonKey192 = []byte{
	0x8e, 0x73, 0xb0, 0xf7, 0xda, 0x0e, 0x64, 0x52, 0xc8, 0x10, 0xf3, 0x2b, 0x80, 0x90, 0x79, 0xe5,
	0x62, 0xf8, 0xea, 0xd2, 0x52, 0x2c, 0x6b, 0x7b,
}

var commonKey256 = []byte{
	0x60, 0x3d, 0xeb, 0x10, 0x15, 0xca, 0x71, 0xbe, 0x2b, 0x73, 0xae, 0xf0, 0x85, 0x7d, 0x77, 0x81,
	0x1f, 0x35, 0x2c, 0x07, 0x3b, 0x61, 0x08, 0xd7, 0x2d, 0x98, 0x10, 0xa3, 0x09, 0x14, 0xdf, 0xf4,
}

func testCipher(t *testing.T, key []byte) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher with %d-bit key: %s", len(key)*8, err)
	}

	encryptedText := make([]byte, aes.BlockSize)
	cipher.Encrypt(encryptedText, commonInput[:aes.BlockSize])

	serializedBlock, err := Serialize(cipher)
	if err != nil {
		t.Fatalf("Failed to serialize cipher with %d-bit key: %s", len(key)*8, err)
	}
	buffer := bytes.NewBuffer(nil)

	err = serializedBlock.Json(buffer)
	if err != nil {
		t.Fatalf("Failed to serialize cipher with %d-bit key: %s", len(key)*8, err)
	}
	bytes, err := buffer.ReadBytes(0)
	if err != nil && err != io.EOF {
		t.Fatalf("Failed to serialize cipher with %d-bit key: %s", len(key)*8, err)
	}
	recreatedBlock, err := NewBlockFromJson(bytes)
	if err != nil {
		t.Fatalf("Failed to deserialize cipher with %d-bit key: %s", len(key)*8, err)
	}
	recoveredCipher, err := Deserialize(recreatedBlock, reflect.TypeOf(cipher))
	if err != nil {
		t.Fatalf("Failed to deserialize cipher with %d-bit key: %s", len(key)*8, err)
	}
	if !reflect.DeepEqual(cipher, recoveredCipher) {
		t.Fatalf("Original and reconstructed cipher are not equal with %d-bit key", len(key)*8)
	}
	encryptedText2 := make([]byte, aes.BlockSize)
	recoveredCipher.Encrypt(encryptedText2, commonInput[:aes.BlockSize])
	if !reflect.DeepEqual(encryptedText, encryptedText2) {
		t.Fatalf("Encrypted text is not equal after reconstruction with %d-bit key", len(key)*8)
	}
}

func TestCiphers(t *testing.T) {
	t.Run("AES128", func(t *testing.T) { testCipher(t, commonKey128) })
	t.Run("AES192", func(t *testing.T) { testCipher(t, commonKey192) })
	t.Run("AES256", func(t *testing.T) { testCipher(t, commonKey256) })
}
