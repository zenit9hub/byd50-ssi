package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
)

// MakeRsaKeys : RSA 개인키/공개키 생성
func MakeRsaKeys(pBits int) (bool, string, string) {
	makePrivateKey, err := rsa.GenerateKey(rand.Reader, pBits)
	if err != nil {
		return false, "", ""
	}
	makePublicKey := &makePrivateKey.PublicKey

	retPrivateKeyString := ExportRSAPrivateKeyAsPEM(makePrivateKey)
	retPublicKeyString, err := ExportRSAPublicKeyAsPEM(makePublicKey)
	if err != nil {
		return false, "", ""
	}

	return true, retPrivateKeyString, retPublicKeyString
}

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Println(err)
	}
	return privkey, &privkey.PublicKey
}
