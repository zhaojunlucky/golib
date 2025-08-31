package security

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestPackDataAndKey(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		key      []byte
		wantSize int
	}{
		{
			name:     "simple data and key",
			data:     []byte("hello world"),
			key:      []byte("secret"),
			wantSize: 4 + len("hello world") + len("secret"), // 4 bytes for key size + data + key
		},
		{
			name:     "empty data",
			data:     []byte{},
			key:      []byte("key123"),
			wantSize: 4 + 0 + len("key123"),
		},
		{
			name:     "empty key",
			data:     []byte("test data"),
			key:      []byte{},
			wantSize: 4 + len("test data") + 0,
		},
		{
			name:     "binary data",
			data:     []byte{0x00, 0x01, 0x02, 0xFF, 0xFE},
			key:      []byte{0xAA, 0xBB, 0xCC},
			wantSize: 4 + 5 + 3,
		},
		{
			name:     "large data",
			data:     make([]byte, 1000),
			key:      make([]byte, 32),
			wantSize: 4 + 1000 + 32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := packDataAndKey(tt.data, tt.key)
			if err != nil {
				t.Errorf("packDataAndKey() error = %v, want nil", err)
				return
			}

			if len(result) != tt.wantSize {
				t.Errorf("packDataAndKey() result size = %d, want %d", len(result), tt.wantSize)
			}

			// Verify the key size is correctly written at the beginning
			buf := bytes.NewBuffer(result)
			var keySize int32
			if err := binary.Read(buf, binary.BigEndian, &keySize); err != nil {
				t.Errorf("Failed to read key size from packed data: %v", err)
			}

			if int(keySize) != len(tt.key) {
				t.Errorf("Packed key size = %d, want %d", keySize, len(tt.key))
			}
		})
	}
}

func TestUnpackDataAndKey(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		key      []byte
		wantErr  bool
	}{
		{
			name:    "simple data and key",
			data:    []byte("hello world"),
			key:     []byte("secret"),
			wantErr: false,
		},
		{
			name:    "empty data",
			data:    []byte{},
			key:     []byte("key123"),
			wantErr: false,
		},
		{
			name:    "empty key",
			data:    []byte("test data"),
			key:     []byte{},
			wantErr: false,
		},
		{
			name:    "binary data",
			data:    []byte{0x00, 0x01, 0x02, 0xFF, 0xFE},
			key:     []byte{0xAA, 0xBB, 0xCC},
			wantErr: false,
		},
		{
			name:    "large data",
			data:    make([]byte, 1000),
			key:     make([]byte, 32),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First pack the data
			packed, err := packDataAndKey(tt.data, tt.key)
			if err != nil {
				t.Errorf("packDataAndKey() error = %v", err)
				return
			}

			// Then unpack it
			unpackedData, unpackedKey, err := unpackDataAndKey(packed)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackDataAndKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !bytes.Equal(unpackedData, tt.data) {
					t.Errorf("unpackDataAndKey() data = %v, want %v", unpackedData, tt.data)
				}

				if !bytes.Equal(unpackedKey, tt.key) {
					t.Errorf("unpackDataAndKey() key = %v, want %v", unpackedKey, tt.key)
				}
			}
		})
	}
}

func TestUnpackDataAndKey_InvalidInput(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		wantErr     bool
		expectPanic bool
	}{
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   []byte{},
			wantErr: true,
		},
		{
			name:    "too short input",
			input:   []byte{0x00, 0x00, 0x00},
			wantErr: true,
		},
		{
			name:        "invalid key size - larger than remaining data",
			input:       []byte{0x00, 0x00, 0x00, 0x10, 0x01, 0x02}, // key size = 16, but only 2 bytes left
			expectPanic: true,
		},
		{
			name:        "negative key size",
			input:       []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x01, 0x02, 0x03, 0x04},
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected function to panic, but it didn't")
					}
				}()
				unpackDataAndKey(tt.input)
			} else {
				_, _, err := unpackDataAndKey(tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("unpackDataAndKey() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestPackUnpackDataAndKey_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		key  []byte
	}{
		{
			name: "text data",
			data: []byte("This is some test data for round trip testing"),
			key:  []byte("encryption_key_123"),
		},
		{
			name: "binary data with nulls",
			data: []byte{0x00, 0x01, 0x00, 0xFF, 0x00, 0xFE, 0x00},
			key:  []byte{0xDE, 0xAD, 0xBE, 0xEF},
		},
		{
			name: "single byte data and key",
			data: []byte{0x42},
			key:  []byte{0x99},
		},
		{
			name: "empty data with key",
			data: []byte{},
			key:  []byte("non_empty_key"),
		},
		{
			name: "data with empty key",
			data: []byte("non_empty_data"),
			key:  []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Pack
			packed, err := packDataAndKey(tt.data, tt.key)
			if err != nil {
				t.Errorf("packDataAndKey() error = %v", err)
				return
			}

			// Unpack
			unpackedData, unpackedKey, err := unpackDataAndKey(packed)
			if err != nil {
				t.Errorf("unpackDataAndKey() error = %v", err)
				return
			}

			// Verify round trip
			if !bytes.Equal(unpackedData, tt.data) {
				t.Errorf("Round trip failed for data: got %v, want %v", unpackedData, tt.data)
			}

			if !bytes.Equal(unpackedKey, tt.key) {
				t.Errorf("Round trip failed for key: got %v, want %v", unpackedKey, tt.key)
			}
		})
	}
}

func TestRandomBytes(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "small size",
			size: 8,
		},
		{
			name: "medium size",
			size: 32,
		},
		{
			name: "large size",
			size: 256,
		},
		{
			name: "single byte",
			size: 1,
		},
		{
			name: "zero size",
			size: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := randomBytes(tt.size)
			if err != nil {
				t.Errorf("randomBytes() error = %v, want nil", err)
				return
			}

			if len(result) != tt.size {
				t.Errorf("randomBytes() size = %d, want %d", len(result), tt.size)
			}
		})
	}
}

func TestRandomBytes_Uniqueness(t *testing.T) {
	size := 32
	iterations := 100

	// Generate multiple random byte arrays and check they're different
	results := make([][]byte, iterations)
	for i := 0; i < iterations; i++ {
		result, err := randomBytes(size)
		if err != nil {
			t.Errorf("randomBytes() error = %v", err)
			return
		}
		results[i] = result
	}

	// Check that all results are different from each other
	for i := 0; i < iterations; i++ {
		for j := i + 1; j < iterations; j++ {
			if bytes.Equal(results[i], results[j]) {
				t.Errorf("randomBytes() generated identical results at indices %d and %d", i, j)
			}
		}
	}
}

func TestRandomBytes_NegativeSize(t *testing.T) {
	// Test with negative size - this should panic as make() cannot create negative-sized slices
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected randomBytes(-1) to panic, but it didn't")
		}
	}()
	
	randomBytes(-1)
}
