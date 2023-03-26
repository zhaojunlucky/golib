package security

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
)

func GenerateECKeyPair(al string) (*ecdsa.PrivateKey, error) {
	switch al {
	case "secp256r1":
		fallthrough
	case "prime256v1":
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "secp224r1":
		return ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "secp384r1":
		return ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "secp521r1":
		return ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return nil, fmt.Errorf("invalid algorithm %s", al)
	}
}

func WritePublicKey(key *ecdsa.PublicKey, writer io.Writer) error {
	x509Encoded, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return err
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509Encoded})
	if _, err := writer.Write(pemEncoded); err != nil {
		return err
	} else {
		return nil
	}
}

func ReadPublicKey(reader io.Reader) (*ecdsa.PublicKey, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	blockPub, _ := pem.Decode(data)
	if blockPub == nil {
		return nil, fmt.Errorf("failed to read public key")
	}
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		return nil, err
	}
	return genericPublicKey.(*ecdsa.PublicKey), nil
}

func WriteECPrivateKey(key *ecdsa.PrivateKey, writer io.Writer) error {
	x509Encoded, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	if _, err := writer.Write(pemEncoded); err != nil {
		return err
	} else {
		return nil
	}
}

func ReadECPrivateKey(reader io.Reader) (*ecdsa.PrivateKey, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to read EC private key")
	}

	x509Encoded := block.Bytes
	key, err := x509.ParsePKCS8PrivateKey(x509Encoded)
	if err != nil {
		return nil, err
	}
	return key.(*ecdsa.PrivateKey), nil
}
