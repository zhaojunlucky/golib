package security

import (
	"crypto/ecdsa"
	"github.com/hotstar/ecies"
)

type ECIESHelper struct {
}

func (eciesHelper *ECIESHelper) EncryptWithPublic(key *ecdsa.PublicKey, data []byte) ([]byte, error) {
	cipher := ecies.NewECIES()

	pubKey := eciesHelper.importPublic(key)
	return cipher.Encrypt(pubKey, data)
}

func (aesHelper *ECIESHelper) importPublic(key *ecdsa.PublicKey) *ecies.PublicKey {
	return &ecies.PublicKey{
		Curve: key.Curve,
		X:     key.X,
		Y:     key.Y,
	}
}

func (eciesHelper *ECIESHelper) importPrivate(key *ecdsa.PrivateKey) *ecies.PrivateKey {
	return &ecies.PrivateKey{
		PublicKey: eciesHelper.importPublic(&key.PublicKey),
		D:         key.D,
	}
}

func (eciesHelper *ECIESHelper) DecryptWithPrivate(key *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	cipher := ecies.NewECIES()
	priKey := eciesHelper.importPrivate(key)
	return cipher.Decrypt(priKey, data)
}
