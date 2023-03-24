package security

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

const (
	IvSize    = 16
	NonceSize = 12
)

type AESHelper struct {
	key     []byte
	TagSize int
}

func NewAESHelper(key []byte) *AESHelper {
	return &AESHelper{key: key, TagSize: 16}
}

func (aesHelper *AESHelper) EncryptCBCRaw(data, iv []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(aesHelper.key)
	if err != nil {
		return nil, nil, err
	}

	padText, err := pkcs7Pad(data, block.BlockSize())
	if err != nil {
		return nil, nil, err
	}

	cipherText := make([]byte, len(padText))

	if iv == nil {
		iv, err = randomBytes(IvSize)
		if err != nil {
			return nil, nil, err
		}
	} else if len(iv) != IvSize {
		return nil, nil, errors.New("invalid IV")
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(cipherText, padText)

	return cipherText, iv, nil
}

func (aesHelper *AESHelper) EncryptCBC(data, iv []byte) ([]byte, error) {
	encryptedData, iv, err := aesHelper.EncryptCBCRaw(data, iv)
	if err != nil {
		return nil, err
	}

	return packDataAndKey(encryptedData, iv)
}

func (aesHelper *AESHelper) DecryptCBC(encryptedData []byte) ([]byte, error) {
	data, iv, err := unpackDataAndKey(encryptedData)
	if err != nil {
		return nil, err
	}

	return aesHelper.DecryptCBCRaw(data, iv)
}

func (aesHelper *AESHelper) DecryptCBCRaw(encryptedData, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesHelper.key)
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(encryptedData, encryptedData)

	cipherText, err := pkcs7Unpad(encryptedData, block.BlockSize())
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}

func (aesHelper *AESHelper) EncryptGCMRaw(data, nonce []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(aesHelper.key)
	if err != nil {
		return nil, nil, err
	}

	if len(nonce) != NonceSize {
		return nil, nil, errors.New("invalid nonce")
	}

	gcm, err := cipher.NewGCMWithTagSize(block, aesHelper.TagSize)
	if err != nil {
		return nil, nil, err
	}

	cipherText := gcm.Seal(nil, nonce, data, nil)

	return cipherText, nonce, nil
}

func (aesHelper *AESHelper) EncryptGCM(data []byte) ([]byte, error) {
	nonce, err := randomBytes(NonceSize)

	if err != nil {
		return nil, err
	}

	encryptedData, nonce, err := aesHelper.EncryptGCMRaw(data, nonce)
	if err != nil {
		return nil, err
	}

	return packDataAndKey(encryptedData, nonce)
}

func (aesHelper *AESHelper) DecryptGCM(encryptedData []byte) ([]byte, error) {
	data, nonce, err := unpackDataAndKey(encryptedData)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(aesHelper.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithTagSize(block, aesHelper.TagSize)
	if err != nil {
		return nil, err
	}

	data, err = gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}
