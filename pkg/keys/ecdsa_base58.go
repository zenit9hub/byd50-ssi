package keys

import (
	"crypto/ecdsa"
	"crypto/x509"
	"log"

	"github.com/btcsuite/btcutil/base58"
)

// ExportECDSAPrivateKeyAsBase58 : Exports a ecdsa.PrivateKey as Base58.
func ExportECDSAPrivateKeyAsBase58(privateKey *ecdsa.PrivateKey) string {
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		log.Printf("error occured: %v", err.Error())
		return ""
	}
	return base58.Encode(privateKeyBytes)
}

// ExportECDSAPublicKeyAsBase58 : Exports a ecdsa.PublicKey as Base58.
func ExportECDSAPublicKeyAsBase58(publicKey *ecdsa.PublicKey) string {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Printf("error occured: %v", err.Error())
		return ""
	}
	return base58.Encode(publicKeyBytes)
}
