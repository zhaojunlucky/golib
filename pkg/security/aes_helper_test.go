package security

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
)

func TestCBC(t *testing.T) {
	aesHelper := NewAESHelper([]byte("passwordpassword"))
	data := []byte("hello world")
	encryptedData, err := aesHelper.EncryptCBC(data, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("encrypted: %s\n", hex.EncodeToString(encryptedData))

	text, err := aesHelper.DecryptCBC(encryptedData)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if !reflect.DeepEqual(data, text) {
		t.Fatalf("CBC encrypt decrypt failed")
	}

}

func TestGCM(t *testing.T) {
	key, _ := hex.DecodeString("24cd27f296351a934855f099c091dc777a8fac258f1fdb7531cd71d7d05f48e0")
	aesHelper := NewAESHelper(key)
	data := []byte("hello world!!!")
	encryptedData, err := aesHelper.EncryptGCM(data)
	if err != nil {
		t.Fatalf(err.Error())
	}

	text, err := aesHelper.DecryptGCM(encryptedData)

	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("encrypted: %s\n", hex.EncodeToString(encryptedData))
	if !reflect.DeepEqual(data, text) {
		t.Fatalf("GCM encrypt decrypt failed")
	}

}

func TestGCMKey(t *testing.T) {
	key, _ := hex.DecodeString("24cd27f296351a934855f099c091dc777a8fac258f1fdb7531cd71d7d05f48e0")
	data := []byte("hello world!!!")
	nonce, _ := hex.DecodeString("ef498a20cd1e4c8fc4712ac3")
	aesHelper := NewAESHelper(key)

	encryptedData, newNonce, err := aesHelper.EncryptGCMRaw(data, nonce)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(hex.EncodeToString(encryptedData))

	if hex.EncodeToString(encryptedData) != "5a0033306a6b97431624d8d9d7733bee8efb6ca189ac7dba25409a536e88" {
		t.Fatalf("GCM encrypted failed")
	}

	encryptedData, _ = packDataAndKey(encryptedData, nonce)

	if !reflect.DeepEqual(newNonce, nonce) {
		t.Fatalf("GCM encrypt decrypt failed")
	}

	text, err := aesHelper.DecryptGCM(encryptedData)

	if err != nil {
		t.Fatalf(err.Error())
	}
	if !reflect.DeepEqual(data, text) {
		t.Fatalf("GCM encrypt decrypt failed")
	}

}
