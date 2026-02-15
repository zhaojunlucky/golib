package security

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"

	"github.com/hotstar/ecies"
)

type ECIESHelper struct {
}

func (eciesHelper *ECIESHelper) EncryptWithPublic(key *ecdsa.PublicKey, data []byte) ([]byte, error) {
	cipher := ecies.NewECIES()

	pubKey, err := eciesHelper.importPublic(key)
	if err != nil {
		return nil, err
	}
	return cipher.Encrypt(pubKey, data)
}

func (eciesHelper *ECIESHelper) importPublic(key *ecdsa.PublicKey) (*ecies.PublicKey, error) {
	pubBytes, err := key.Bytes()
	if err != nil {
		return nil, err
	}
	x, y := elliptic.Unmarshal(key.Curve, pubBytes)
	if x == nil || y == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return &ecies.PublicKey{Curve: key.Curve, X: x, Y: y}, nil
}

func (eciesHelper *ECIESHelper) importPrivate(key *ecdsa.PrivateKey) (*ecies.PrivateKey, error) {
	pubKey, err := eciesHelper.importPublic(&key.PublicKey)
	if err != nil {
		return nil, err
	}
	privBytes, err := key.Bytes()
	if err != nil {
		return nil, err
	}
	return &ecies.PrivateKey{
		PublicKey: pubKey,
		D:         new(big.Int).SetBytes(privBytes),
	}, nil
}

func (eciesHelper *ECIESHelper) DecryptWithPrivate(key *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	cipher := ecies.NewECIES()
	priKey, err := eciesHelper.importPrivate(key)
	if err != nil {
		return nil, err
	}
	return cipher.Decrypt(priKey, data)
}
