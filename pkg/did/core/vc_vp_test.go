package core_test

import (
	"byd50-ssi/pkg/did/core"
	byd50_jwt "byd50-ssi/pkg/did/core/byd50-jwt"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"testing"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/golang-jwt/jwt"
)

func TestVerifyVcAndVp(t *testing.T) {
	pvKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	pbBytes, err := x509.MarshalPKIXPublicKey(&pvKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	pbKeyBase58 := base58.Encode(pbBytes)

	getPbKey := func(_ string, _ string) string {
		return pbKeyBase58
	}

	credSub := map[string]interface{}{
		"name": "tester",
	}
	standardClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "did:byd50:test",
	}

	vcJwt := core.CreateVc("did:byd50:test", "TestCredential", credSub, standardClaims, pvKey)
	ok, err := core.VerifyVc(vcJwt, getPbKey)
	if !ok || err != nil {
		t.Fatalf("verify vc failed: %v", err)
	}

	vpJwt := core.CreateVp("did:byd50:test", "TestPresentation", []string{vcJwt}, standardClaims, pvKey)
	ok, _, err = core.VerifyVp(vpJwt, getPbKey)
	if !ok || err != nil {
		t.Fatalf("verify vp failed: %v", err)
	}

	ok, claims, err := core.GetMapClaims(vpJwt, getPbKey)
	if !ok || err != nil {
		t.Fatalf("get map claims failed: %v", err)
	}
	if _, err := core.ClaimsGetExp(claims); err != nil {
		t.Fatalf("claims exp failed: %v", err)
	}
	if _, err := core.ClaimsGetIat(claims); err != nil {
		t.Fatalf("claims iat failed: %v", err)
	}
}

func TestValidateClaims(t *testing.T) {
	now := time.Now()

	vcClaims := byd50_jwt.VcClaims{
		Nonce: "n",
		Vc:    map[string]interface{}{"type": []string{"VerifiableCredential"}},
		StandardClaims: jwt.StandardClaims{
			Issuer:    "issuer",
			ExpiresAt: now.Add(time.Minute).Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}
	if err := core.ValidateVcClaims(vcClaims); err != nil {
		t.Fatalf("vc claims should be valid: %v", err)
	}
	vcClaims.ExpiresAt = now.Add(-time.Minute).Unix()
	if err := core.ValidateVcClaims(vcClaims); err == nil {
		t.Fatal("expected vc claims expiry error")
	}

	vpClaims := byd50_jwt.VpClaims{
		Nonce: "n",
		Vp:    map[string]interface{}{"type": []string{"VerifiablePresentation"}},
		StandardClaims: jwt.StandardClaims{
			Issuer:    "issuer",
			ExpiresAt: now.Add(time.Minute).Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}
	if err := core.ValidateVpClaims(vpClaims); err != nil {
		t.Fatalf("vp claims should be valid: %v", err)
	}
	vpClaims.Issuer = ""
	if err := core.ValidateVpClaims(vpClaims); err == nil {
		t.Fatal("expected vp claims issuer error")
	}
}
