package security

import (
	"bytes"
	"errors"
)

var (
	InvalidBlockSize = errors.New("invalid block size")

	InvalidPKCS7Data = errors.New("invalid PKCS7 data (empty or not padded)")

	InvalidPKCS7Padding = errors.New("invalid padding on input")
)

func pkcs7Pad(b []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, InvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, InvalidPKCS7Data
	}
	n := blockSize - (len(b) % blockSize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

func pkcs7Unpad(b []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, InvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, InvalidPKCS7Data
	}
	if len(b)%blockSize != 0 {
		return nil, InvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, InvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, InvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}
