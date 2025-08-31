package security

import (
	"bytes"
	"testing"
)

func TestPkcs7Pad_ValidInput(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		blockSize int
		expected  []byte
	}{
		{
			name:      "pad 1 byte to 8-byte block",
			input:     []byte("hello"),
			blockSize: 8,
			expected:  []byte("hello\x03\x03\x03"),
		},
		{
			name:      "pad to 16-byte block",
			input:     []byte("hello world"),
			blockSize: 16,
			expected:  []byte("hello world\x05\x05\x05\x05\x05"),
		},
		{
			name:      "pad single byte",
			input:     []byte("a"),
			blockSize: 4,
			expected:  []byte("a\x03\x03\x03"),
		},
		{
			name:      "pad exactly one block size",
			input:     []byte("12345678"),
			blockSize: 8,
			expected:  []byte("12345678\x08\x08\x08\x08\x08\x08\x08\x08"),
		},
		{
			name:      "pad with block size 1",
			input:     []byte("test"),
			blockSize: 1,
			expected:  []byte("test\x01"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pkcs7Pad(tt.input, tt.blockSize)
			if err != nil {
				t.Errorf("pkcs7Pad() error = %v, want nil", err)
				return
			}
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("pkcs7Pad() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPkcs7Pad_InvalidBlockSize(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		blockSize int
	}{
		{
			name:      "zero block size",
			input:     []byte("test"),
			blockSize: 0,
		},
		{
			name:      "negative block size",
			input:     []byte("test"),
			blockSize: -1,
		},
		{
			name:      "very negative block size",
			input:     []byte("test"),
			blockSize: -100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pkcs7Pad(tt.input, tt.blockSize)
			if err != ErrInvalidBlockSize {
				t.Errorf("pkcs7Pad() error = %v, want %v", err, ErrInvalidBlockSize)
			}
			if result != nil {
				t.Errorf("pkcs7Pad() result = %v, want nil", result)
			}
		})
	}
}

func TestPkcs7Pad_InvalidData(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		blockSize int
	}{
		{
			name:      "nil input",
			input:     nil,
			blockSize: 8,
		},
		{
			name:      "empty input",
			input:     []byte{},
			blockSize: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pkcs7Pad(tt.input, tt.blockSize)
			if err != ErrInvalidPKCS7Data {
				t.Errorf("pkcs7Pad() error = %v, want %v", err, ErrInvalidPKCS7Data)
			}
			if result != nil {
				t.Errorf("pkcs7Pad() result = %v, want nil", result)
			}
		})
	}
}

func TestPkcs7Unpad_ValidInput(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		blockSize int
		expected  []byte
	}{
		{
			name:      "unpad 8-byte block",
			input:     []byte("hello\x03\x03\x03"),
			blockSize: 8,
			expected:  []byte("hello"),
		},
		{
			name:      "unpad 16-byte block",
			input:     []byte("hello world\x05\x05\x05\x05\x05"),
			blockSize: 16,
			expected:  []byte("hello world"),
		},
		{
			name:      "unpad single byte padding",
			input:     []byte("test\x01"),
			blockSize: 5,
			expected:  []byte("test"),
		},
		{
			name:      "unpad full block padding",
			input:     []byte("12345678\x08\x08\x08\x08\x08\x08\x08\x08"),
			blockSize: 8,
			expected:  []byte("12345678"),
		},
		{
			name:      "unpad with large padding",
			input:     []byte("a\x07\x07\x07\x07\x07\x07\x07"),
			blockSize: 8,
			expected:  []byte("a"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pkcs7Unpad(tt.input, tt.blockSize)
			if err != nil {
				t.Errorf("pkcs7Unpad() error = %v, want nil", err)
				return
			}
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("pkcs7Unpad() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPkcs7Unpad_InvalidBlockSize(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		blockSize int
	}{
		{
			name:      "zero block size",
			input:     []byte("test\x01"),
			blockSize: 0,
		},
		{
			name:      "negative block size",
			input:     []byte("test\x01"),
			blockSize: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pkcs7Unpad(tt.input, tt.blockSize)
			if err != ErrInvalidBlockSize {
				t.Errorf("pkcs7Unpad() error = %v, want %v", err, ErrInvalidBlockSize)
			}
			if result != nil {
				t.Errorf("pkcs7Unpad() result = %v, want nil", result)
			}
		})
	}
}

func TestPkcs7Unpad_InvalidData(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		blockSize int
		wantErr   error
	}{
		{
			name:      "nil input",
			input:     nil,
			blockSize: 8,
			wantErr:   ErrInvalidPKCS7Data,
		},
		{
			name:      "empty input",
			input:     []byte{},
			blockSize: 8,
			wantErr:   ErrInvalidPKCS7Data,
		},
		{
			name:      "input not multiple of block size",
			input:     []byte("hello"),
			blockSize: 8,
			wantErr:   ErrInvalidPKCS7Padding,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pkcs7Unpad(tt.input, tt.blockSize)
			if err != tt.wantErr {
				t.Errorf("pkcs7Unpad() error = %v, want %v", err, tt.wantErr)
			}
			if result != nil {
				t.Errorf("pkcs7Unpad() result = %v, want nil", result)
			}
		})
	}
}

func TestPkcs7Unpad_InvalidPadding(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		blockSize int
	}{
		{
			name:      "zero padding byte",
			input:     []byte("hello\x00\x00\x00"),
			blockSize: 8,
		},
		{
			name:      "padding byte larger than input",
			input:     []byte("hi\x08\x08"),
			blockSize: 4,
		},
		{
			name:      "inconsistent padding bytes",
			input:     []byte("test\x03\x02\x03\x03"),
			blockSize: 8,
		},
		{
			name:      "padding byte too large for block",
			input:     []byte("test\x09\x09\x09\x09"),
			blockSize: 8,
		},
		{
			name:      "single byte with wrong padding",
			input:     []byte("a\x02\x02\x03"),
			blockSize: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pkcs7Unpad(tt.input, tt.blockSize)
			if err != ErrInvalidPKCS7Padding {
				t.Errorf("pkcs7Unpad() error = %v, want %v", err, ErrInvalidPKCS7Padding)
			}
			if result != nil {
				t.Errorf("pkcs7Unpad() result = %v, want nil", result)
			}
		})
	}
}

func TestPkcs7PadUnpad_RoundTrip(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		blockSize int
	}{
		{
			name:      "short string",
			input:     []byte("hello"),
			blockSize: 8,
		},
		{
			name:      "long string",
			input:     []byte("this is a longer test string for padding"),
			blockSize: 16,
		},
		{
			name:      "single byte",
			input:     []byte("x"),
			blockSize: 4,
		},
		{
			name:      "exact block size",
			input:     []byte("12345678"),
			blockSize: 8,
		},
		{
			name:      "binary data",
			input:     []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD},
			blockSize: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Pad the data
			padded, err := pkcs7Pad(tt.input, tt.blockSize)
			if err != nil {
				t.Errorf("pkcs7Pad() error = %v", err)
				return
			}

			// Unpad the data
			unpadded, err := pkcs7Unpad(padded, tt.blockSize)
			if err != nil {
				t.Errorf("pkcs7Unpad() error = %v", err)
				return
			}

			// Verify round trip
			if !bytes.Equal(unpadded, tt.input) {
				t.Errorf("Round trip failed: original = %v, got = %v", tt.input, unpadded)
			}
		})
	}
}
