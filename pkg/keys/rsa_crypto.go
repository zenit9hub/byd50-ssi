package keys

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"log"
)

// RsaEnc : RSA 암호화
func RsaEnc(pPublicKey string, pMessage string, pLabel string) (bool, string) {
	byteMessage := []byte(pMessage)
	// byteLabel := []byte(pLabel)
	// hash := sha256.New()

	publicKey, err := ParseRSAPublicKeyFromPEM(pPublicKey)
	if err != nil {
		return false, ""
	}

	// cipherText, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, byteMessage, byteLabel)
	// if err != nil {
	// 	return false, ""
	// }

	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, byteMessage)
	if err != nil {
		return false, ""
	}

	return true, base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(cipherText)))
}

// RsaDec : RSA 복호화
func RsaDec(pPrivateKey string, pCipherMessage string, pLabel string) (bool, string) {
	byteMessage, err := base64.StdEncoding.DecodeString(pCipherMessage)
	if err != nil {
		return false, ""
	}

	byteHexDecode, err := hex.DecodeString(string(byteMessage))
	if err != nil {
		return false, ""
	}

	// byteLabel := []byte(pLabel)
	// hash := sha256.New()

	privateKey, err := ParseRSAPrivateKeyFromPEM(pPrivateKey)
	if err != nil {
		return false, ""
	}

	// decryptedBytes, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, byteMessage, byteLabel)
	// if err != nil {
	// 	return false, ""
	// }

	// decryptedBytes, err := privateKey.Decrypt(nil, byteMessage, &rsa.OAEPOptions{Hash: crypto.SHA256})
	// if err != nil {
	// 	panic(err)
	// }

	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, byteHexDecode)
	if err != nil {
		panic(err)
	}

	return true, string(decryptedBytes)
}

// RsaSign : RSA Sign
func RsaSign(pPrivateKey string, pCipherMessage string, pLabel string) (bool, string) {
	// byteLabel := []byte(pLabel)
	hash := sha256.New()
	hash.Write([]byte(pCipherMessage))
	digest := hash.Sum(nil)

	privateKey, err := ParseRSAPrivateKeyFromPEM(pPrivateKey)
	if err != nil {
		return false, ""
	}

	signBytes, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest)
	if err != nil {
		panic(err)
	}

	return true, base64.StdEncoding.EncodeToString(signBytes)
}

// RsaVerify : RSA Verify
func RsaVerify(pPublicKey string, pCipherMessage string, pSign string) bool {
	// byteCipherMessage, err := base64.StdEncoding.DecodeString(pCipherMessage)
	// if err != nil {
	// 	return false
	// }

	// byteSign, err := base64.StdEncoding.DecodeString(pSign)
	// if err != nil {
	// 	return false
	// }

	hash := sha256.New()
	hash.Write([]byte(pCipherMessage))
	digest := hash.Sum(nil)

	publicKey, err := ParseRSAPublicKeyFromPEM(pPublicKey)
	if err != nil {
		return false
	}

	//err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, digest, []byte(pSign))
	byteSign, _ := base64.StdEncoding.DecodeString(pSign)
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, digest, byteSign)
	if err != nil {
		return false
	}

	return true
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		log.Println(err)
	}
	return ciphertext
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		log.Println(err)
	}
	return plaintext
}
