package security

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestAESHelper_EncryptWithPrivate(t *testing.T) {
	ecies := ECIESHelper{}

	priKey, err := GenerateECKeyPair("secp256r1")
	if err != nil {
		t.Fatal(err)
	}
	writer := new(bytes.Buffer)

	WriteECPrivateKey(priKey, writer)
	fmt.Printf("private key: %s\n", writer.String())

	message := "hello world!!!"

	encrypted, err := ecies.EncryptWithPublic(&priKey.PublicKey, []byte(message))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("encrypted data: %s\n", hex.EncodeToString(encrypted))

	decrypted, _ := ecies.DecryptWithPrivate(priKey, encrypted)
	decryptedText := string(decrypted)
	if message != decryptedText {
		t.Fatal("failed to verify encryption and decryption")
	}
}
