package security

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
)

type AES interface {
	encryptCBCRaw(data, iv []byte) ([]byte, []byte, error)
	encryptCBC(data, iv []byte) ([]byte, error)
	decryptCBC(encryptedData, iv []byte) ([]byte, error)

	encryptGCMRaw(data, nonce []byte) ([]byte, []byte, error)
	encryptGCM(data, nonce []byte) ([]byte, error)
	decryptGCM(encryptedData []byte) ([]byte, error)
}

func packDataAndKey(data, key []byte) ([]byte, error) {
	keySize := len(key)

	buf := &bytes.Buffer{}
	buf.Grow(len(data) + keySize + 4)

	err := binary.Write(buf, binary.BigEndian, int32(keySize))
	if err != nil {
		return nil, err
	}

	if _, err := buf.Write(data); err != nil {
		return nil, err
	}

	if _, err := buf.Write(key); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func unpackDataAndKey(data []byte) ([]byte, []byte, error) {
	size := len(data)
	buf := bytes.NewBuffer(data)
	var keySize int32
	if err := binary.Read(buf, binary.BigEndian, &keySize); err != nil {
		return nil, nil, err
	}

	dataSize := size - 4 - int(keySize)
	packData := make([]byte, dataSize)
	keyData := make([]byte, keySize)

	if _, err := buf.Read(packData); err != nil {
		return nil, nil, err
	}

	if _, err := buf.Read(keyData); err != nil {
		return nil, nil, err
	}
	return packData, keyData, nil
}

func randomBytes(size int) ([]byte, error) {
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		return nil, err
	}
	return data, nil
}
