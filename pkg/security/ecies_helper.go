package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

type ECIESHelper struct {
}

// EncryptWithPublic encrypts data using ECIES with the recipient's public key.
// Wire format: [ephemeral_public_key || ciphertext || gcm_tag]
func (eciesHelper *ECIESHelper) EncryptWithPublic(key *ecdsa.PublicKey, data []byte) ([]byte, error) {
	// Convert ECDSA public key to ECDH
	ecdhPub, err := ecdsaPublicToECDH(key)
	if err != nil {
		return nil, err
	}

	// Generate ephemeral key pair
	ephemeralPriv, err := ecdhPub.Curve().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	// Perform ECDH to get shared secret
	sharedSecret, err := ephemeralPriv.ECDH(ecdhPub)
	if err != nil {
		return nil, err
	}

	// Derive encryption key using HKDF-SHA256
	kdf := hkdf.New(sha256.New, sharedSecret, nil, nil)
	encKey := make([]byte, 32) // AES-256
	if _, err := io.ReadFull(kdf, encKey); err != nil {
		return nil, err
	}

	// Encrypt with AES-GCM
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	// Prepend ephemeral public key
	ephemeralPubBytes := ephemeralPriv.PublicKey().Bytes()
	result := make([]byte, len(ephemeralPubBytes)+len(ciphertext))
	copy(result, ephemeralPubBytes)
	copy(result[len(ephemeralPubBytes):], ciphertext)

	return result, nil
}

// DecryptWithPrivate decrypts ECIES-encrypted data using the recipient's private key.
func (eciesHelper *ECIESHelper) DecryptWithPrivate(key *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	// Convert ECDSA private key to ECDH
	ecdhPriv, err := ecdsaPrivateToECDH(key)
	if err != nil {
		return nil, err
	}

	// Determine ephemeral public key size based on curve
	var ephemeralPubSize int
	switch ecdhPriv.Curve() {
	case ecdh.P256():
		ephemeralPubSize = 65 // 1 + 32 + 32
	case ecdh.P384():
		ephemeralPubSize = 97 // 1 + 48 + 48
	case ecdh.P521():
		ephemeralPubSize = 133 // 1 + 66 + 66
	default:
		return nil, fmt.Errorf("unsupported curve")
	}

	if len(data) < ephemeralPubSize+12+16 { // min: pubkey + nonce + tag
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract ephemeral public key
	ephemeralPubBytes := data[:ephemeralPubSize]
	ephemeralPub, err := ecdhPriv.Curve().NewPublicKey(ephemeralPubBytes)
	if err != nil {
		return nil, err
	}

	// Perform ECDH to get shared secret
	sharedSecret, err := ecdhPriv.ECDH(ephemeralPub)
	if err != nil {
		return nil, err
	}

	// Derive decryption key using HKDF-SHA256
	kdf := hkdf.New(sha256.New, sharedSecret, nil, nil)
	decKey := make([]byte, 32) // AES-256
	if _, err := io.ReadFull(kdf, decKey); err != nil {
		return nil, err
	}

	// Decrypt with AES-GCM
	block, err := aes.NewCipher(decKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := data[ephemeralPubSize:]
	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// ecdsaPublicToECDH converts an ECDSA public key to ECDH public key
func ecdsaPublicToECDH(key *ecdsa.PublicKey) (*ecdh.PublicKey, error) {
	pubBytes, err := key.Bytes()
	if err != nil {
		return nil, err
	}

	var curve ecdh.Curve
	switch key.Curve.Params().Name {
	case "P-224":
		return nil, fmt.Errorf("P-224 not supported by crypto/ecdh")
	case "P-256":
		curve = ecdh.P256()
	case "P-384":
		curve = ecdh.P384()
	case "P-521":
		curve = ecdh.P521()
	default:
		return nil, fmt.Errorf("unsupported curve: %s", key.Curve.Params().Name)
	}

	return curve.NewPublicKey(pubBytes)
}

// ecdsaPrivateToECDH converts an ECDSA private key to ECDH private key
func ecdsaPrivateToECDH(key *ecdsa.PrivateKey) (*ecdh.PrivateKey, error) {
	privBytes, err := key.Bytes()
	if err != nil {
		return nil, err
	}

	var curve ecdh.Curve
	switch key.Curve.Params().Name {
	case "P-224":
		return nil, fmt.Errorf("P-224 not supported by crypto/ecdh")
	case "P-256":
		curve = ecdh.P256()
	case "P-384":
		curve = ecdh.P384()
	case "P-521":
		curve = ecdh.P521()
	default:
		return nil, fmt.Errorf("unsupported curve: %s", key.Curve.Params().Name)
	}

	return curve.NewPrivateKey(privBytes)
}
