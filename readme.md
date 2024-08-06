# CipherBlock Serializer

## Overview
The CipherBlock Serializer is a Go library designed to serialize and deserialize cryptographic cipher blocks (`cipher.Block`) into byte slices (`[]byte`). This functionality is particularly useful for embedding cryptographic cipher configurations directly into code, facilitating secure, hardcoded cryptographic setups.

## Features
- **Serialize Cipher Blocks**: Convert any `cipher.Block` into a byte slice, capturing the exact state and configuration of the cipher.
- **Deserialize Cipher Blocks**: Reconstruct a `cipher.Block` from a byte slice, ensuring that the cryptographic properties are retained.
- **Support for Multiple Key Sizes**: Compatible with AES-128, AES-192, and AES-256 key sizes.

## Installation
To install CipherBlock Serializer, use the following `go get` command:

```bash
go get -u github.com/RafOSS-br/go-cipher-block-serializer
