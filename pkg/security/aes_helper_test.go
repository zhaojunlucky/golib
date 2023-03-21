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
	encryptedData, err := aesHelper.encryptCBC(data, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("encrypted: %s\n", hex.EncodeToString(encryptedData))

	text, err := aesHelper.decryptCBC(encryptedData)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if !reflect.DeepEqual(data, text) {
		t.Fatalf("CBC encrypt decrypt failed")
	}

}

func TestGCM(t *testing.T) {
	aesHelper := NewAESHelper([]byte("passwordpassword"))
	data := []byte("hello world")
	encryptedData, err := aesHelper.encryptGCM(data, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	text, err := aesHelper.decryptGCM(encryptedData)

	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Printf("encrypted: %s\n", hex.EncodeToString(encryptedData))
	if !reflect.DeepEqual(data, text) {
		t.Fatalf("GCM encrypt decrypt failed")
	}

}
