package keys

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// ExportECDSAPrivateKeyAsPEM : Exports a ecdsa.PrivateKey as PEM.
func ExportECDSAPrivateKeyAsPEM(privateKey *ecdsa.PrivateKey) (string, error) {
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)
	return string(privateKeyPEM), nil
}

// ParseECDSAPrivateKeyFromPEM : Parses a ecdsa.PrivateKey from PEM.
func ParseECDSAPrivateKeyFromPEM(privateKeyPEM string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// ExportECDSAPublicKeyAsPEM : Exports a ecdsa.PublicKey as PEM.
func ExportECDSAPublicKeyAsPEM(publicKey *ecdsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)
	return string(publicKeyPEM), nil
}

// ParseECDSAPublicKeyFromPEM : Parses a ecdsa.PublicKey from PEM.
func ParseECDSAPublicKeyFromPEM(publicKeyPEM string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch pub := publicKey.(type) {
	case *ecdsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("key type is not ecdsa.PublicKey")
	}
}
