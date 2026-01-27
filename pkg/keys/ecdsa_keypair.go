package keys

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

// MakeEcdsaKeys : ECDSA 개인키/공개키 생성(PEM 반환)
func MakeEcdsaKeys() (bool, string, string) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return false, "", ""
	}
	publicKey := &privateKey.PublicKey

	privateKeyPem, err := ExportECDSAPrivateKeyAsPEM(privateKey)
	if err != nil {
		return false, "", ""
	}
	publicKeyPem, err := ExportECDSAPublicKeyAsPEM(publicKey)
	if err != nil {
		return false, "", ""
	}

	return true, privateKeyPem, publicKeyPem
}

// GenerateECDSAKeyPair generates a new ECDSA key pair.
func GenerateECDSAKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}
