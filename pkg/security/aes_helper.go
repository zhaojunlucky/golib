package security

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

const (
	IV_SIZE    = 16
	NONCE_SIZE = 12
)

type AESHelper struct {
	key []byte
}

func NewAESHelper(key []byte) *AESHelper {
	return &AESHelper{key: key}
}

func (aesHelper *AESHelper) encryptCBCRaw(data, iv []byte) ([]byte, []byte, error) {
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
		iv, err = randomBytes(IV_SIZE)
		if err != nil {
			return nil, nil, err
		}
	} else if len(iv) != IV_SIZE {
		return nil, nil, errors.New("invalid IV")
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(cipherText, padText)

	return cipherText, iv, nil
}

func (aesHelper *AESHelper) encryptCBC(data, iv []byte) ([]byte, error) {
	encryptedData, iv, err := aesHelper.encryptCBCRaw(data, iv)
	if err != nil {
		return nil, err
	}

	return packDataAndKey(encryptedData, iv)
}

func (aesHelper *AESHelper) decryptCBC(encryptedData []byte) ([]byte, error) {
	data, iv, err := unpackDataAndKey(encryptedData)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(aesHelper.key)
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(data, data)

	cipherText, err := pkcs7Unpad(data, block.BlockSize())
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}

func (aesHelper *AESHelper) encryptGCMRaw(data, nonce []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(aesHelper.key)
	if err != nil {
		return nil, nil, err
	}

	if nonce == nil {
		nonce, err = randomBytes(NONCE_SIZE)

		if err != nil {
			return nil, nil, err
		}
	} else if len(nonce) != NONCE_SIZE {
		return nil, nil, errors.New("invalid nonce")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	cipherText := gcm.Seal(nil, nonce, data, nil)

	return cipherText, nonce, nil
}

func (aesHelper *AESHelper) encryptGCM(data, nonce []byte) ([]byte, error) {
	encryptedData, nonce, err := aesHelper.encryptGCMRaw(data, nonce)
	if err != nil {
		return nil, err
	}

	return packDataAndKey(encryptedData, nonce)
}

func (aesHelper *AESHelper) decryptGCM(encryptedData []byte) ([]byte, error) {
	data, nonce, err := unpackDataAndKey(encryptedData)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(aesHelper.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	cipherText, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}
