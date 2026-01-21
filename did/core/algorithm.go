package core

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"github.com/btcsuite/btcutil/base58"
	"log"
	mathrand "math/rand"
)

func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	return hex.EncodeToString(bytes), err
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()-=_+")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[mathrand.Intn(len(letters))]
	}
	return string(s)
}

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

//func Contain(target interface{}, list interface{}) (bool, int) {
//	if reflect.TypeOf(list).Kind() == reflect.Slice || reflect.TypeOf(list).Kind() == reflect.Array {
//		listvalue := reflect.ValueOf(list)
//		for i := 0; i < listvalue.Len(); i++ {
//			if target == listvalue.Index(i).Interface() {
//				return true, i
//			}
//		}
//	}
//	if reflect.TypeOf(target).Kind() == reflect.String && reflect.TypeOf(list).Kind() == reflect.String {
//		return strings.Contains(list.(string), target.(string)), strings.Index(list.(string), target.(string))
//	}
//	return false, -1
//}

func PbKeyEncrypt(pbKeyBase58, plainText string) string {
	pbKey, _ := x509.ParsePKCS1PublicKey(base58.Decode(pbKeyBase58))

	_, cipherText := rsaEncrypt(plainText, pbKey)

	cipherTextBase64 := base64.StdEncoding.EncodeToString(cipherText)
	log.Printf("  -- plainText: %v", plainText)
	log.Printf("  -- cipherText(%v)", len(cipherText))
	log.Printf("  == cipherTextBase64(%v)", len(cipherTextBase64))

	return cipherTextBase64
}

func PvKeyDecrypt(ciphertextBase64, pvKeyBase58 string) string {
	ciphertext, _ := base64.StdEncoding.DecodeString(ciphertextBase64)
	pvKey, _ := x509.ParsePKCS1PrivateKey(base58.Decode(pvKeyBase58))

	_, plaintext := rsaDecrypt(ciphertext, pvKey)
	return string(plaintext)
}

func PvKeySign(pvKeyBase58 string, pCipherMessage string, pLabel string) (bool, string) {
	result := false
	var signBytes []byte
	pvKey, err := x509.ParsePKCS1PrivateKey(base58.Decode(pvKeyBase58))
	if err == nil {
		result, signBytes = rsaSign(pCipherMessage, pvKey)
	}
	return result, base64.StdEncoding.EncodeToString(signBytes)
}

func PbKeyVerify(pbKeyBase58 string, pCipherMessage string, pSign string) bool {
	pbKey, _ := x509.ParsePKCS1PublicKey(base58.Decode(pbKeyBase58))

	result := rsaVerify(pCipherMessage, pSign, pbKey)

	return result
}

func rsaSign(pCipherMessage string, pvKey *rsa.PrivateKey) (bool, []byte) {
	hash := sha256.New()
	hash.Write([]byte(pCipherMessage))
	digest := hash.Sum(nil)
	result := false

	signBytes, err := rsa.SignPKCS1v15(rand.Reader, pvKey, crypto.SHA256, digest)
	if err == nil {
		result = true
	}

	return result, signBytes
}

func rsaVerify(pCipherMessage, pSign string, pbKey *rsa.PublicKey) bool {
	hash := sha256.New()
	hash.Write([]byte(pCipherMessage))
	digest := hash.Sum(nil)
	result := false

	byteSign, _ := base64.StdEncoding.DecodeString(pSign)
	err := rsa.VerifyPKCS1v15(pbKey, crypto.SHA256, digest, byteSign)
	if err == nil {
		result = true
	}

	return result
}

func rsaEncrypt(plainText string, pbKey *rsa.PublicKey) (bool, []byte) {
	hash := sha512.New()
	cipherText, err := rsa.EncryptOAEP(hash, rand.Reader, pbKey, []byte(plainText), nil)
	result := false

	if err == nil {
		result = true
	}

	return result, cipherText
}

func rsaDecrypt(ciphertext []byte, priv *rsa.PrivateKey) (bool, []byte) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	result := false

	if err == nil {
		result = true
	}
	return result, plaintext
}
