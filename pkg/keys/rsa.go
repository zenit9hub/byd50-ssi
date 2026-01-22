package keys

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
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

// ExportRSAPrivateKeyAsPEM : Exports a rsa.PrivateKey as PEM.
func ExportRSAPrivateKeyAsPEM(privateKey *rsa.PrivateKey) string {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)
	return string(privateKeyPEM)
}

// ParseRSAPrivateKeyFromPEM : Parses a rsa.PrivateKey from PEM.
func ParseRSAPrivateKeyFromPEM(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// ExportRSAPublicKeyAsPEM : Exports a rsa.PublicKey as PEM.
func ExportRSAPublicKeyAsPEM(publicKey *rsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)

	return string(publicKeyPEM), nil
}

// ParseRSAPublicKeyFromPEM : Parses a rsa.PublicKey from PEM.
func ParseRSAPublicKeyFromPEM(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := publicKey.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("key type is not rsa.PublicKey")
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// ExportPublicKeyAsPemStr : 공개키 PEM 문자열 변환
func ExportPublicKeyAsPemStr(pubkey *rsa.PublicKey) string {
	pubkeyPem := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(pubkey)}))
	return pubkeyPem
}

// ExportPrivateKeyAsPemStr : 개인키 PEM 문자열 변환
func ExportPrivateKeyAsPemStr(privatekey *rsa.PrivateKey) string {
	privatekeyPem := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privatekey)}))
	return privatekeyPem
}

// ExportMsgAsPemStr :
func ExportMsgAsPemStr(msg []byte) string {
	msgPem := string(pem.EncodeToMemory(&pem.Block{Type: "MESSAGE", Bytes: msg}))
	return msgPem
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// ExportRsaPrivateKeyAsPemStr :
func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	return string(privkeyPem)
}

// ParseRsaPrivateKeyFromPemStr :
func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

// ExportRsaPublicKeyAsPemStr :
func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	pubkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)

	return string(pubkeyPem), nil
}

// ParseRsaPublicKeyFromPemStr :
func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("Key type is not RSA")
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Println(err)
	}
	return privkey, &privkey.PublicKey
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		log.Println(err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Println(err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		log.Println(err)
	}
	return key
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Println(err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		log.Println(err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Println("not ok")
	}
	return key
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

// bits := 2048
// bobPrivateKey, _ := rsa.GenerateKey(rand.Reader, bits)

// bobPublicKey := &bobPrivateKey.PublicKey

// fmt.Printf("%s\n", utility.ExportPrivateKeyAsPemStr(bobPrivateKey))

// fmt.Printf("%s\n", utility.ExportPublicKeyAsPemStr(bobPublicKey))

// message := []byte("test")
// label := []byte("")
// hash := sha256.New()

// ciphertext, _ := rsa.EncryptOAEP(hash, rand.Reader, bobPublicKey, message, label)

// fmt.Printf("%s\n", utility.ExportMsgAsPemStr(ciphertext))

// plainText, _ := rsa.DecryptOAEP(hash, rand.Reader, bobPrivateKey, ciphertext, label)

// fmt.Printf("RSA decrypted to [%s]", plainText)

// // Create the keys
// priv, pub := GenerateRsaKeyPair()

// // Export the keys to pem string
// priv_pem := ExportRsaPrivateKeyAsPemStr(priv)
// pub_pem, _ := ExportRsaPublicKeyAsPemStr(pub)

// // Import the keys from pem string
// priv_parsed, _ := ParseRsaPrivateKeyFromPemStr(priv_pem)
// pub_parsed, _ := ParseRsaPublicKeyFromPemStr(pub_pem)

// // Export the newly imported keys
// priv_parsed_pem := ExportRsaPrivateKeyAsPemStr(priv_parsed)
// pub_parsed_pem, _ := ExportRsaPublicKeyAsPemStr(pub_parsed)

// fmt.Println(priv_parsed_pem)
// fmt.Println(pub_parsed_pem)

// // Check that the exported/imported keys match the original keys
// if priv_pem != priv_parsed_pem || pub_pem != pub_parsed_pem {
// 		fmt.Println("Failure: Export and Import did not result in same Keys")
// } else {
// 		fmt.Println("Success")
// }
