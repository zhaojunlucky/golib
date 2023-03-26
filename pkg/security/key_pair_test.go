package security

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestReadPublicKey(t *testing.T) {
	keyStr := `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE6BzTH3rLmg7f/ErVIfS/bkx8DglC
tU4K4u52kLxUvdkMkE59ktO+Q+iVL3aOlZOG4iMcMscdCt3G601RmdmP5A==
-----END PUBLIC KEY-----`

	pub, err := ReadPublicKey(strings.NewReader(keyStr))
	if err != nil {
		t.Fatalf("Failed to read public key: %v", err)
	}

	writer := new(bytes.Buffer)

	err = WritePublicKey(pub, writer)
	if err != nil {
		t.Fatalf("failed to write public key %v", err)
	}

	pubStr := strings.TrimSpace(writer.String())

	if !reflect.DeepEqual(keyStr, pubStr) {
		t.Fatalf("Failed to verify public key")
	}
}

func TestWriteECPrivateKey(t *testing.T) {
	key, err := GenerateECKeyPair("secp256r1")
	if err != nil {
		t.Fatal(err)
	}

	writer := new(bytes.Buffer)

	WriteECPrivateKey(key, writer)
	fmt.Println(writer.String())
}

func TestReadECPrivateKey(t *testing.T) {
	keyStr := `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgrF/+jV+eT8gukETa
PmnVU7evsRBcxVvR9MXIILkv/iWgCgYIKoZIzj0DAQehRANCAARSfH94I5c40duO
kJWr3E224SDmpcJUHjj9tylsdl/OetC2Rrn/99lEQ4g5WskBeFqiKRWmDbO7kgsc
Whn99/HP
-----END PRIVATE KEY-----`

	pri, err := ReadECPrivateKey(strings.NewReader(keyStr))
	if err != nil {
		t.Fatalf("Failed to read private key: %v", err)
	}

	writer := new(bytes.Buffer)

	err = WriteECPrivateKey(pri, writer)
	if err != nil {
		t.Fatalf("failed to write EC private key %v", err)
	}
	pubStr := strings.TrimSpace(writer.String())

	newPri, err := ReadECPrivateKey(strings.NewReader(pubStr))
	if err != nil {
		t.Fatal(err)
	}
	if !newPri.Equal(pri) {
		t.Fatal("failed to verify private key")
	}
}
