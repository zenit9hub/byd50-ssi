package keys

import (
	"crypto/rsa"
	"crypto/x509"
	"log"

	"github.com/btcsuite/btcutil/base58"
)

// ExportRSAPrivateKeyAsBase58 : Exports a rsa.PrivateKey as Base58.
func ExportRSAPrivateKeyAsBase58(privateKey *rsa.PrivateKey) string {
	if privateKey == nil {
		return ""
	}
	return base58.Encode(x509.MarshalPKCS1PrivateKey(privateKey))
}

// ExportRSAPublicKeyAsBase58 : Exports a rsa.PublicKey as Base58.
func ExportRSAPublicKeyAsBase58(publicKey *rsa.PublicKey) string {
	if publicKey == nil {
		return ""
	}
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	if len(publicKeyBytes) == 0 {
		log.Printf("empty rsa public key bytes")
		return ""
	}
	return base58.Encode(publicKeyBytes)
}
