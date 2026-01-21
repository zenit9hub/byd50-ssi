package utility

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"
)

func TestRsaEncDec(t *testing.T) {
	_, privPem, pubPem := MakeRsaKeys(2048)
	if pubPem == "" || privPem == "" {
		t.Fatal("failed to generate keys")
	}

	ok, cipher := RsaEnc(pubPem, "hello", "")
	if !ok {
		t.Fatal("encrypt failed")
	}

	ok, plain := RsaDec(privPem, cipher, "")
	if !ok {
		t.Fatal("decrypt failed")
	}
	if plain != "hello" {
		t.Fatalf("unexpected plaintext: %s", plain)
	}

	if ok, _ := RsaEnc("bad-key", "hello", ""); ok {
		t.Fatal("expected encrypt failure")
	}
	if ok, _ := RsaDec("bad-key", "bad", ""); ok {
		t.Fatal("expected decrypt failure")
	}
	if ok, _ := RsaDec("bad-key", cipher, ""); ok {
		t.Fatal("expected decrypt failure with bad private key")
	}
	if ok, _ := RsaDec(privPem, "***", ""); ok {
		t.Fatal("expected decrypt base64 failure")
	}
	invalidHex := base64.StdEncoding.EncodeToString([]byte("zz"))
	if ok, _ := RsaDec(privPem, invalidHex, ""); ok {
		t.Fatal("expected decrypt hex failure")
	}
}

func TestRsaSignVerify(t *testing.T) {
	_, privPem, pubPem := MakeRsaKeys(2048)
	ok, sign := RsaSign(privPem, "payload", "")
	if !ok {
		t.Fatal("sign failed")
	}
	if !RsaVerify(pubPem, "payload", sign) {
		t.Fatal("verify failed")
	}

	if RsaVerify("bad-key", "payload", sign) {
		t.Fatal("expected verify failure")
	}
	if ok, _ := RsaSign("bad-key", "payload", ""); ok {
		t.Fatal("expected sign failure")
	}
}

func TestKeyConversions(t *testing.T) {
	priv, pub := GenerateKeyPair(2048)
	privBytes := PrivateKeyToBytes(priv)
	pubBytes := PublicKeyToBytes(pub)

	if len(privBytes) == 0 || len(pubBytes) == 0 {
		t.Fatal("key bytes empty")
	}

	privParsed := BytesToPrivateKey(privBytes)
	pubParsed := BytesToPublicKey(pubBytes)
	if privParsed == nil || pubParsed == nil {
		t.Fatal("parsed keys are nil")
	}

	msg := []byte("secret")
	cipher := EncryptWithPublicKey(msg, pubParsed)
	plain := DecryptWithPrivateKey(cipher, privParsed)
	if !bytes.Equal(msg, plain) {
		t.Fatal("encrypt/decrypt mismatch")
	}

	encryptedPriv, err := x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(priv), []byte("pw"), x509.PEMCipherAES256)
	if err != nil {
		t.Fatal(err)
	}
	_ = BytesToPrivateKey(pem.EncodeToMemory(encryptedPriv))

	encryptedPub, err := x509.EncryptPEMBlock(rand.Reader, "RSA PUBLIC KEY", pubBytes, []byte("pw"), x509.PEMCipherAES256)
	if err != nil {
		t.Fatal(err)
	}
	_ = BytesToPublicKey(pem.EncodeToMemory(encryptedPub))
}

func TestPemHelpers(t *testing.T) {
	priv, pub := GenerateKeyPair(2048)

	if ExportPublicKeyAsPemStr(pub) == "" {
		t.Fatal("ExportPublicKeyAsPemStr empty")
	}
	if ExportPrivateKeyAsPemStr(priv) == "" {
		t.Fatal("ExportPrivateKeyAsPemStr empty")
	}
	if ExportMsgAsPemStr([]byte("msg")) == "" {
		t.Fatal("ExportMsgAsPemStr empty")
	}

	privPem := ExportRsaPrivateKeyAsPemStr(priv)
	if privPem == "" {
		t.Fatal("ExportRsaPrivateKeyAsPemStr empty")
	}
	if _, err := ParseRsaPrivateKeyFromPemStr(privPem); err != nil {
		t.Fatal(err)
	}

	pubPem, err := ExportRsaPublicKeyAsPemStr(pub)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseRsaPublicKeyFromPemStr(pubPem); err != nil {
		t.Fatal(err)
	}

	privPem2 := ExportRSAPrivateKeyAsPEM(priv)
	if _, err := ParseRSAPrivateKeyFromPEM(privPem2); err != nil {
		t.Fatal(err)
	}
	pubPem2, err := ExportRSAPublicKeyAsPEM(pub)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseRSAPublicKeyFromPEM(pubPem2); err != nil {
		t.Fatal(err)
	}

	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	ecdsaPubBytes, err := x509.MarshalPKIXPublicKey(&ecdsaKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	ecdsaPubPem := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: ecdsaPubBytes})
	if _, err := ParseRSAPublicKeyFromPEM(string(ecdsaPubPem)); err == nil {
		t.Fatal("expected rsa public key parse error")
	}
	if _, err := ParseRsaPublicKeyFromPemStr(string(ecdsaPubPem)); err == nil {
		t.Fatal("expected rsa public key parse error (pem str)")
	}
	ecdsaPrivBytes, err := x509.MarshalECPrivateKey(ecdsaKey)
	if err != nil {
		t.Fatal(err)
	}
	ecdsaPrivPem := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: ecdsaPrivBytes})
	if _, err := ParseRSAPrivateKeyFromPEM(string(ecdsaPrivPem)); err == nil {
		t.Fatal("expected rsa private key parse error")
	}

	if _, err := ParseRsaPrivateKeyFromPemStr("bad-key"); err == nil {
		t.Fatal("expected parse private key error")
	}
	if _, err := ParseRsaPublicKeyFromPemStr("bad-key"); err == nil {
		t.Fatal("expected parse public key error")
	}
}
