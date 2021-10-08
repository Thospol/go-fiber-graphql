package utils

import (
	"crypto/ecdsa"
	"os"

	"github.com/dgrijalva/jwt-go"
)

var (
	// SignKey is private key
	SignKey = &ecdsa.PrivateKey{}
	// VerifyKey is public key
	VerifyKey = &ecdsa.PublicKey{}
)

// ReadECDSAKey read private key && public key
func ReadECDSAKey(privateKey, publicKey string) error {
	privateKeyByte, err := os.ReadFile(privateKey)
	if err != nil {
		return err
	}

	publicKeyByte, err := os.ReadFile(publicKey)
	if err != nil {
		return err
	}

	SignKey, err = jwt.ParseECPrivateKeyFromPEM(privateKeyByte)
	if err != nil {
		return err
	}

	VerifyKey, err = jwt.ParseECPublicKeyFromPEM(publicKeyByte)
	if err != nil {
		return err
	}

	return nil
}
